package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/accrual"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/balance"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/order"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type AppContainer struct {
	Config                *config.Config
	Logger                logger.Logger
	DB                    *sql.DB
	AuthContainer         *auth.AuthContainer
	OrderContainer        *order.OrderContainer
	BalanceContainer      *balance.BalanceContainer
	AccrualOrderProcessor *accrual.OrderProcessor
	AccrualHTTPClient     *accrual.HTTPAccrualClient
}

func NewAppContainer(ctx context.Context, config *config.Config) (*AppContainer, error) {
	logger := logger.New("info")

	defer logger.Sync()

	pool, err := pgxpool.New(ctx, config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	userRepo := repository.NewDBUserRepository(db)
	orderRepo := repository.NewDBOrderRepository(db)
	balanceRepo := repository.NewDBBalanceRepository(db)

	accrualClient := accrual.NewHTTPAccrualClient(fmt.Sprintf("%s/api/orders/", config.AccrualAddress) + "%s")
	accrualOrderProcessor := accrual.NewOrderProcessor(20, orderRepo, logger, accrualClient)

	authContainer := auth.NewContainer(config, logger, userRepo)
	orderContainer := order.NewContainer(config, logger, authContainer.TokenService, orderRepo)
	balanceContainer := balance.NewContainer(config, logger, authContainer.TokenService, balanceRepo)

	return &AppContainer{
		Config:                config,
		Logger:                logger,
		DB:                    db,
		AuthContainer:         authContainer,
		OrderContainer:        orderContainer,
		BalanceContainer:      balanceContainer,
		AccrualOrderProcessor: accrualOrderProcessor,
		AccrualHTTPClient:     accrualClient,
	}, nil
}
