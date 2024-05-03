package auth

import (
	"context"
	"errors"

	proto "github.com/sodiqit/gophermart/gen/proto"
	"github.com/sodiqit/gophermart/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	proto.UnimplementedAuthServer
	logger      logger.Logger
	authService AuthService
}

func (s *AuthServer) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	var response proto.LoginResponse

	logger := s.logger.With("op", proto.Auth_Login_FullMethodName)

	result, err := s.authService.Login(ctx, in.Email, in.Password)

	if err != nil {
		return nil, mapLoginServiceError(err, logger)
	}

	response.Token = result

	return &response, nil
}

func mapLoginServiceError(err error, logger logger.Logger) error {
	code := codes.Internal
	msg := "Internal server error"

	if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrIncorrectPassword) {
		code = codes.Unauthenticated
		msg = err.Error()
	}

	return status.Error(code, msg)
}

func NewAuthServer(logger logger.Logger, authService AuthService) *AuthServer {
	return &AuthServer{
		logger:      logger,
		authService: authService,
	}
}
