package balance

import (
	"context"
	"errors"
	"fmt"

	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

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
	op := "balanceService.withdraw"

	balance, err := s.balanceRepo.GetBalanceWithWithdrawals(ctx, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if balance.Current-sum < 0 {
		return ErrInsufficientFunds
	}

	_, err = s.balanceRepo.CreateWithdraw(ctx, userID, orderID, sum)

	return err
}

func (s *SimpleBalanceService) GetWithdrawals(ctx context.Context, userID int) ([]dtos.Withdraw, error) {
	return s.balanceRepo.GetWithdrawalsByUser(ctx, userID)
}

func NewService(balanceRepo repository.BalanceRepository) *SimpleBalanceService {
	return &SimpleBalanceService{
		balanceRepo: balanceRepo,
	}
}
