package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/order"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

func RunServer(config *config.Config) error {
	logger := logger.New("info")

	defer logger.Sync()

	pool, err := pgxpool.New(context.Background(), config.DatabaseDSN)
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)

	err = db.Ping()

	if err != nil {
		return err
	}

	userRepo := repository.NewDBUserRepository(db)
	orderRepo := repository.NewDBOrderRepository(db)

	authContainer := auth.NewContainer(config, logger, userRepo)
	orderContainer := order.NewContainer(config, logger, authContainer.TokenService, orderRepo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/user", authContainer.Controller.Route())
	r.Mount("/api/user/orders", orderContainer.Controller.Route())

	logger.Infow("start server", "address", config.Address, "config", config)
	return http.ListenAndServe(config.Address, r)
}
