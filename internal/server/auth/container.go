package auth

import (
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type AuthContainer struct {
	TokenService      TokenService
	SimpleAuthService AuthService
	Controller        *AuthController
	GRPCServer        *AuthServer
}

func NewContainer(config *config.Config, logger logger.Logger, userRepo repository.UserRepository) *AuthContainer {
	tokenService := NewJWTTokenService(config.JWTSecretKey, config.JWTTimeExp)
	authService := NewSimpleAuthService(tokenService, userRepo)
	authController := NewController(logger, authService)
	authServer := NewAuthServer(logger, authService)

	return &AuthContainer{
		TokenService:      tokenService,
		SimpleAuthService: authService,
		Controller:        authController,
		GRPCServer:        authServer,
	}
}
