package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/sodiqit/gophermart/gen/gophermart_db/public/model"
	"github.com/sodiqit/gophermart/gen/gophermart_db/public/table"
	"github.com/sodiqit/gophermart/internal/server/dtos"
)

var ErrOrderNotFound = errors.New("order not found")

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type OrderRepository interface {
	Create(ctx context.Context, userID int, orderNumber string, status string) (string, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (dtos.Order, error)
	GetListByUser(ctx context.Context, userID int) ([]dtos.Order, error)
	GetOrdersForProcessing(ctx context.Context, pool int64) ([]string, error)
	UpdateOrder(ctx context.Context, orderID string, status string, accrual *float64) error
}

type DBOrderRepository struct {
	db *sql.DB
}

func (r *DBOrderRepository) Create(ctx context.Context, userID int, orderNumber string, status string) (string, error) {
	op := "orderRepo.create"

	stmt := table.Orders.INSERT(table.Orders.ID, table.Orders.UserID, table.Orders.Status).
		VALUES(orderNumber, userID, status).
		RETURNING(table.Orders.ID)

	var dest model.Orders

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return dest.ID, nil
}

func (r *DBOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (dtos.Order, error) {
	op := "orderRepo.findByOrderNumber"

	stmt := table.Orders.SELECT(table.Orders.ID, table.Orders.UserID, table.Orders.Accrual, table.Orders.Status, table.Orders.CreatedAt, table.Orders.UpdatedAt).WHERE(table.Orders.ID.EQ(postgres.String(orderNumber)))

	var dest model.Orders

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil && errors.Is(qrm.ErrNoRows, err) {
		return dtos.Order{}, fmt.Errorf("%s: %w", op, ErrOrderNotFound)
	}

	if err != nil {
		return dtos.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	return mapOrderEntityToDto(dest), nil
}

func (r *DBOrderRepository) GetListByUser(ctx context.Context, userID int) ([]dtos.Order, error) {
	op := "orderRepo.getListByUser"

	stmt := table.Orders.SELECT(table.Orders.ID, table.Orders.UserID, table.Orders.Accrual, table.Orders.Status, table.Orders.CreatedAt, table.Orders.UpdatedAt).WHERE(table.Orders.UserID.EQ(postgres.Int(int64(userID))))

	var dest []model.Orders

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return make([]dtos.Order, 0), fmt.Errorf("%s: %w", op, err)
	}

	result := make([]dtos.Order, len(dest))

	for i, entity := range dest {
		result[i] = mapOrderEntityToDto(entity)
	}

	return result, nil
}

func (r *DBOrderRepository) GetOrdersForProcessing(ctx context.Context, pool int64) ([]string, error) {
	op := "orderRepo.getOrdersForProcessing"

	stmt := table.Orders.SELECT(table.Orders.ID).
		WHERE(
			table.Orders.Status.IN(postgres.String(OrderStatusNew), postgres.String(OrderStatusProcessing)),
		).
		LIMIT(pool)

	var dest []model.Orders

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return make([]string, 0), fmt.Errorf("%s: %w", op, err)
	}

	result := make([]string, len(dest))

	for i, entity := range dest {
		result[i] = entity.ID
	}

	return result, nil
}

func (r *DBOrderRepository) UpdateOrder(ctx context.Context, orderID string, status string, accrual *float64) error {
	stmt := table.Orders.UPDATE(table.Orders.Status, table.Orders.Accrual).SET(status, accrual).WHERE(table.Orders.ID.EQ(postgres.String(orderID)))

	_, err := stmt.ExecContext(ctx, r.db)

	return err
}

func mapOrderEntityToDto(entity model.Orders) dtos.Order {
	return dtos.Order{
		ID:        entity.ID,
		UserID:    int(entity.UserID),
		Accrual:   entity.Accrual,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

}

var _ OrderRepository = (*DBOrderRepository)(nil)

func NewDBOrderRepository(db *sql.DB) *DBOrderRepository {
	return &DBOrderRepository{db: db}
}
