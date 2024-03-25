package auth

import (
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/config"
)

type AuthContainer struct {
	TokenService      TokenService
	SimpleAuthService AuthService
	Controller        *AuthController
}

func NewContainer(config *config.Config, logger logger.Logger) *AuthContainer {
	tokenService := NewJWTTokenService(config.JWTSecretKey)
	authService := NewSimpleAuthService(tokenService)
	authController := NewController(logger, authService, tokenService)

	return &AuthContainer{
		TokenService:      tokenService,
		SimpleAuthService: authService,
		Controller:        authController,
	}
}
