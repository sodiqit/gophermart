package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/config"
)

func RunServer(config *config.Config) error {
	logger := logger.New("info")

	defer logger.Sync()

	authController := auth.NewController()

	r := chi.NewRouter()
	r.Mount("/", authController.Route())

	logger.Infow("start server", "address", config.Address, "config", config)
	return http.ListenAndServe(config.Address, r)
}
