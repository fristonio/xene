package agent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/fristonio/xene/pkg/apiserver/client"
	"github.com/fristonio/xene/pkg/apiserver/client/registry"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	types "github.com/fristonio/xene/pkg/types/v1alpha1"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	log "github.com/sirupsen/logrus"
)

// JoinAPIServer joins the agent to the pool of servers registered on the API
// server.
func (s *Server) JoinAPIServer(apiServerAddr, agentName, agentAddr, authToken string) error {
	addr, err := url.Parse(apiServerAddr)
	if err != nil {
		return fmt.Errorf("error while parsing Agent address: %s", err)
	}

	client := client.New(httptransport.New(addr.Host, "", nil), strfmt.Default)
	bearerTokenAuth := httptransport.BearerToken(authToken)

	res, err := client.Registry.GetAPIV1RegistryAgent(
		registry.NewGetAPIV1RegistryAgentParams().WithName(&agentName), bearerTokenAuth)
	if err != nil {
		return fmt.Errorf("error getting current agent information: %s", err)
	}

	a := types.Agent{}
	if res.Payload.Item.Value != "" {
		err = json.Unmarshal([]byte(res.Payload.Item.Value), &a)
		if err != nil {
			return fmt.Errorf("error unmarshaling agent get response: %s", err)
		}
	}

	agentManifest := types.Agent{
		TypeMeta: types.TypeMeta{
			Kind:       types.AgentKind,
			APIVersion: "v1alpha1",
		},
		Metadata: types.Metadata{
			ObjectMeta: types.ObjectMeta{
				Name:        agentName,
				Description: "Worker agents for xene",
			},
		},
		Spec: types.AgentSpec{
			Address:          agentAddr,
			Insecure:         option.Config.Agent.Insecure,
			LogServerEnabled: option.Config.Agent.LocalLogServer,
			LogServerPort:    option.Config.Agent.LogServerPort,
		},
	}

	if !option.Config.Agent.Insecure {
		// First create secrets corresponding to the agent.
		err = s.createXeneSecrets(addr.Host, agentName, authToken)
		if err != nil {
			return err
		}

		agentManifest.Spec.ClientKeySecret = fmt.Sprintf("%s-client-key", agentName)
		agentManifest.Spec.ClientCertSecret = fmt.Sprintf("%s-client-cert", agentName)
		agentManifest.Spec.RootCASecret = fmt.Sprintf("%s-root-ca-chain", agentName)
		agentManifest.Spec.ServerName = option.Config.Agent.ServerName
	}

	log.Infof("%v", a)
	log.Infof("%v", agentManifest)
	if a.DeepEquals(&agentManifest) {
		log.Infof("agent information is already up to date in apiserver")
		return nil
	}

	log.Infof("creating agent object for apiserver")

	ag, err := json.Marshal(agentManifest)
	if err != nil {
		return fmt.Errorf("error while marshaling agent manifest: %s", err)
	}

	resp, err := client.Registry.
		PostAPIV1RegistryAgent(registry.NewPostAPIV1RegistryAgentParams().WithAgent(string(ag)), bearerTokenAuth)
	if err != nil {
		return err
	}

	log.Infof("created %s agent: %s", agentName, resp.Payload.Message)
	return nil
}

// createXeneSecrets creates all the required xene secrets in xene registry
// These secrets then can be in turn used by any agent GRPC client to manage
// agent.
func (s *Server) createXeneSecrets(apiServerAddr, agentName, authToken string) error {
	if option.Config.Agent.Insecure {
		return nil
	}

	log.Infof("creating xene agent client certificates which can be used by apiserver")
	clientKey, err := ioutil.ReadFile(option.Config.Agent.ClientKeyFile)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %s", option.Config.Agent.ClientKeyFile, err)
	}

	clientCert, err := ioutil.ReadFile(option.Config.Agent.ClientCertFile)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %s", option.Config.Agent.ClientCertFile, err)
	}

	rootCA, err := ioutil.ReadFile(option.Config.Agent.RootCACert)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %s", option.Config.Agent.RootCACert, err)
	}

	client := client.New(httptransport.New(apiServerAddr, "", nil), strfmt.Default)
	bearerTokenAuth := httptransport.BearerToken(authToken)

	var (
		secretData []byte
		secret     = &v1alpha1.Secret{
			TypeMeta: v1alpha1.TypeMeta{
				Kind:       v1alpha1.SecretKind,
				APIVersion: "v1alpha1",
			},
			Metadata: v1alpha1.Metadata{
				ObjectMeta: types.ObjectMeta{
					Description: "Secrets for xene agent.",
				},
			},
			Spec: v1alpha1.SecretSpec{
				Type: defaults.SecretTypeCertificates,
			},
		}
	)

	secret.Spec.Content = base64.StdEncoding.EncodeToString([]byte(clientKey))
	secret.Metadata.ObjectMeta.Name = fmt.Sprintf("%s-client-key", agentName)
	secretData, err = json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("error while marshaling secret manifest: %s", err)
	}
	resp, err := client.Registry.
		PostAPIV1RegistrySecret(registry.NewPostAPIV1RegistrySecretParams().WithSecret(string(secretData)), bearerTokenAuth)
	if err != nil {
		return err
	}
	log.Infof("created %s secret: %s", secret.Metadata.GetName(), resp.Payload.Message)

	secret.Spec.Content = base64.StdEncoding.EncodeToString([]byte(clientCert))
	secret.Metadata.ObjectMeta.Name = fmt.Sprintf("%s-client-cert", agentName)
	secretData, err = json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("error while marshaling secret manifest: %s", err)
	}
	resp, err = client.Registry.
		PostAPIV1RegistrySecret(registry.NewPostAPIV1RegistrySecretParams().WithSecret(string(secretData)), bearerTokenAuth)
	if err != nil {
		return err
	}
	log.Infof("created %s secret: %s", secret.Metadata.GetName(), resp.Payload.Message)

	secret.Spec.Content = base64.StdEncoding.EncodeToString([]byte(rootCA))
	secret.Metadata.ObjectMeta.Name = fmt.Sprintf("%s-root-ca-chain", agentName)
	secretData, err = json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("error while marshaling secret manifest: %s", err)
	}
	resp, err = client.Registry.
		PostAPIV1RegistrySecret(registry.NewPostAPIV1RegistrySecretParams().WithSecret(string(secretData)), bearerTokenAuth)
	if err != nil {
		return err
	}
	log.Infof("created %s secret: %s", secret.Metadata.GetName(), resp.Payload.Message)

	return nil
}
