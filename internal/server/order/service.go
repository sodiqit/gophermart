package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/sodiqit/gophermart/internal/server/repository"
)

type OrderService interface {
	Upload(ctx context.Context, userID int, orderNumber string) error
}

var ErrUserAlreadyUploadOrder = errors.New("user already upload this order")
var ErrOrderAlreadyUploadByAnotherUser = errors.New("another user already upload this order")

type SimpleOrderService struct {
	orderRepo repository.OrderRepository
}

func (s *SimpleOrderService) Upload(ctx context.Context, userID int, orderNumber string) error {
	op := "orderService.upload"

	order, err := s.orderRepo.FindByOrderNumber(ctx, orderNumber)

	if err == nil && order.UserID == userID {
		return fmt.Errorf("%s: %w", op, ErrUserAlreadyUploadOrder)
	}

	if err == nil && order.UserID != userID {
		return fmt.Errorf("%s: %w", op, ErrOrderAlreadyUploadByAnotherUser)
	}

	if err != nil && !errors.Is(err, repository.ErrOrderNotFound) {
		return err
	}

	_, err = s.orderRepo.Create(ctx, userID, orderNumber, repository.OrderStatusNew)

	return err
}

func NewSimpleOrderService(orderRepo repository.OrderRepository) *SimpleOrderService {
	return &SimpleOrderService{
		orderRepo: orderRepo,
	}
}
