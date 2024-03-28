package auth

import (
	"fmt"
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

	r.Post("/register", c.handleRegister)

	return r
}

func (c *AuthController) handleRegister(w http.ResponseWriter, r *http.Request) {
	op := "authController.handleRegister"

	token, err := c.authService.Register(r.Context(), "test", "test")

	if err != nil {
		c.logger.Errorw("op", op, "err", err.Error())
		http.Error(w, "", http.StatusInternalServerError) //TODO: map errors from service to http status codes
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

func NewController(logger logger.Logger, authService AuthService, tokenService TokenService) *AuthController {
	return &AuthController{
		logger,
		authService,
		tokenService,
	}
}
