package balance

import (
	"context"

	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type BalanceService interface {
	GetTotalBalance(ctx context.Context, userID int) (dtos.Balance, error)
	Withdraw(ctx context.Context, userID int, orderID string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]dtos.Withdraw, error)
}

type SimpleBalanceService struct {
	balanceRepo repository.BalanceRepository
}

func (s *SimpleBalanceService) GetTotalBalance(ctx context.Context, userID int) (dtos.Balance, error) {
	return s.balanceRepo.GetBalanceWithWithdrawals(ctx, userID)
}

func (s *SimpleBalanceService) Withdraw(ctx context.Context, userID int, orderID string, sum float64) error {
	return nil //TODO: implement
}

func (s *SimpleBalanceService) GetWithdrawals(ctx context.Context, userID int) ([]dtos.Withdraw, error) {
	return nil, nil //TODO: implement
}

func NewService(balanceRepo repository.BalanceRepository) *SimpleBalanceService {
	return &SimpleBalanceService{
		balanceRepo: balanceRepo,
	}
}

