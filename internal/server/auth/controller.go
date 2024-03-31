package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/utils"
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
	r.Post("/login", c.handleLogin)

	return r
}

func (c *AuthController) handleRegister(w http.ResponseWriter, r *http.Request) {
	op := "authController.handleRegister"

	logger := c.logger.With("op", op)

	var dto RegisterRequestDTO

	err := utils.ValidateJSONBody(r.Context(), r.Body, &dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := c.authService.Register(r.Context(), dto.Username, dto.Password)

	if err != nil {
		mapRegisterErrorToHttpError(w, err, logger, dto)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (c *AuthController) handleLogin(w http.ResponseWriter, r *http.Request) {
	op := "authController.handleLogin"

	logger := c.logger.With("op", op)

	var dto LoginRequestDTO

	err := utils.ValidateJSONBody(r.Context(), r.Body, &dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(r.Context(), dto.Username, dto.Password)

	if err != nil {
		mapLoginErrorToHttpError(w, err, logger, dto)
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

func mapRegisterErrorToHttpError(w http.ResponseWriter, err error, logger logger.Logger, dto RegisterRequestDTO) {
	if errors.Is(err, ErrUserAlreadyExist) {
		logger.Infow("", "username", dto.Username, "err", err.Error())
		http.Error(w, "", http.StatusConflict)
		return
	}

	logger.Errorw("", "err", err.Error(), "username", dto.Username)
	http.Error(w, "", http.StatusInternalServerError)
}

func mapLoginErrorToHttpError(w http.ResponseWriter, err error, logger logger.Logger, dto LoginRequestDTO) {
	if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrIncorrectPassword) {
		logger.Infow("", "username", dto.Username, "err", err.Error())
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	logger.Errorw("", "err", err.Error(), "username", dto.Username)
	http.Error(w, "", http.StatusInternalServerError)
}
