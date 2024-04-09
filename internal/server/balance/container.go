package balance

import (
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type BalanceContainer struct {
	Controller *BalanceController
	Service    BalanceService
}

func NewContainer(config *config.Config, logger logger.Logger, tokenService auth.TokenService, balanceRepo repository.BalanceRepository) *BalanceContainer {
	service := NewService(balanceRepo)
	controller := NewController(logger, tokenService, service)

	return &BalanceContainer{
		Controller: controller,
		Service:    service,
	}
}
