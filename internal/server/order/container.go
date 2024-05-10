package order

import (
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type OrderContainer struct {
	Controller *OrderController
	GRPCServer *OrderServer
}

func NewContainer(config *config.Config, logger logger.Logger, tokenService auth.TokenService, orderRepo repository.OrderRepository) *OrderContainer {
	orderService := NewSimpleOrderService(orderRepo)
	orderController := NewController(logger, tokenService, orderService)
	orderServer := NewOrderServer(logger, orderService)

	return &OrderContainer{
		Controller: orderController,
		GRPCServer: orderServer,
	}
}
