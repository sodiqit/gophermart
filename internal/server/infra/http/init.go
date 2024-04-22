package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/accrual"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/balance"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/order"
	"github.com/sodiqit/gophermart/internal/server/repository"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	balanceRepo := repository.NewDBBalanceRepository(db)

	accrualClient := accrual.NewHTTPAccrualClient(fmt.Sprintf("%s/api/orders/", config.AccrualAddress) + "%s")
	accrualOrderProcessor := accrual.NewOrderProcessor(20, orderRepo, logger, accrualClient)

	authContainer := auth.NewContainer(config, logger, userRepo)
	orderContainer := order.NewContainer(config, logger, authContainer.TokenService, orderRepo)
	balanceContainer := balance.NewContainer(config, logger, authContainer.TokenService, balanceRepo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))
	r.Mount("/debug", middleware.Profiler())
	r.Mount("/api/user", authContainer.Controller.Route())
	r.Mount("/api/user/orders", orderContainer.Controller.Route())
	balanceContainer.Controller.Connect(r, "/api/")

	go accrualOrderProcessor.Run(context.TODO()) //TODO: use shutdown context

	logger.Infow("start server", "address", config.Address, "config", config)
	return http.ListenAndServe(config.Address, r)
}
