package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fristonio/xene/pkg/auth"
	types "github.com/fristonio/xene/pkg/types/v1alpha1"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// JWTVerficationMiddleware is the middleware handler to authenticate requests to agents
// GRPC server.
func (s *Server) JWTVerficationMiddleware(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := s.authProvider.GetServerClaimsFromToken(token)
	if err != nil {
		return nil, err
	}

	log.Debugf("Server claims: %v", claims)
	return handler(ctx, req)
}

// JoinAPIServer joins the agent to the pool of servers registered on the API
// server.
func (s *Server) JoinAPIServer(apiServerAddr, agentName, agentAddr, authToken string) error {
	agentAuthToken, err := s.authProvider.NewServerAuthToken(apiServerAddr, "apiserver")
	if err != nil {
		return fmt.Errorf("error while generating auth token: %s", err)
	}
	agentManifest := types.Agent{
		TypeMeta: types.TypeMeta{
			Kind:       "agent",
			APIVersion: "v1alpha1",
		},
		Metadata: types.Metadata{
			ObjectMeta: types.ObjectMeta{
				Name:        agentName,
				Description: "Worker agents for xene",
			},
		},
		Spec: types.AgentSpec{
			Address:   agentAddr,
			AuthToken: agentAuthToken,
		},
	}

	ag, err := json.Marshal(agentManifest)
	if err != nil {
		return fmt.Errorf("error while marshaling agent manifest: %s", err)
	}

	// Prepare request to be sent to the API server for registring the agent.
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/registry/agent", apiServerAddr),
		strings.NewReader(url.Values{
			"agent": []string{string(ag)},
		}.Encode()))
	if err != nil {
		return fmt.Errorf("error while preparing request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while making request to apiserver: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid request: %s", string(respBody))
	}

	log.Debug("response from API server while registering: %v", string(respBody))
	return nil
}
