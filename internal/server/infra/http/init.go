package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sodiqit/gophermart/internal/server/infra"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var srv http.Server

func RunServer(ctx context.Context, deps *infra.AppContainer) error {
	config := deps.Config
	logger := deps.Logger
	authContainer := deps.AuthContainer
	orderContainer := deps.OrderContainer
	balanceContainer := deps.BalanceContainer
	accrualOrderProcessor := deps.AccrualOrderProcessor

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
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.Write([]byte("pong"))
	})
	balanceContainer.Controller.Connect(r, "/api/")

	go accrualOrderProcessor.Run(ctx)

	logger.Infow("start HTTP server", "address", config.Address, "config", config)
	srv = http.Server{Addr: config.Address, Handler: r}
	return srv.ListenAndServe()
}

func StopServer(ctx context.Context) error {
	fmt.Println("Gracefully shutdown HTTP server")
	return srv.Shutdown(ctx)
}
