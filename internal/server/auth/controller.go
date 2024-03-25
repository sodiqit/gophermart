package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sodiqit/gophermart/internal/logger"
)

type AuthController struct {
	logger       logger.Logger
	authService  AuthService
	tokenService TokenService
}

func (c *AuthController) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(JWTAuth(c.tokenService))

	r.Post("/register", c.handleRegister)

	return r
}

func (c *AuthController) handleRegister(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewController(logger logger.Logger, authService AuthService, tokenService TokenService) *AuthController {
	return &AuthController{
		logger,
		authService,
		tokenService,
	}
}
