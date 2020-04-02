package agent

import (
	"context"

	"github.com/fristonio/xene/pkg/auth"
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
