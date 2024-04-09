package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/sodiqit/gophermart/gen/gophermart_db/public/model"
	"github.com/sodiqit/gophermart/gen/gophermart_db/public/table"
	"github.com/sodiqit/gophermart/internal/server/dtos"
)

type BalanceRepository interface {
	GetBalanceWithWithdrawals(ctx context.Context, userID int) (dtos.Balance, error)
	CreateWithdraw(ctx context.Context, userID int, orderID string, sum float64) (int, error)
	GetWithdrawalsByUser(ctx context.Context, userID int) ([]dtos.Withdraw, error)
}

type DBBalanceRepository struct {
	db *sql.DB
}

func (r *DBBalanceRepository) GetBalanceWithWithdrawals(ctx context.Context, userID int) (dtos.Balance, error) {
	op := "balanceRepo.getBalanceWithWithdrawals"

	query := `
		WITH accrued AS (
			SELECT 
				user_id, 
				SUM(accrual) AS total_accrued
			FROM 
				orders
			WHERE 
				status = $2
			GROUP BY 
				user_id
		), withdrawn AS (
			SELECT 
				user_id, 
				SUM(amount) AS total_withdrawn
			FROM 
				withdraws
			GROUP BY 
				user_id
		)
		SELECT 
			u.id AS user_id, 
			COALESCE(a.total_accrued, 0) - COALESCE(w.total_withdrawn, 0) AS current_balance,
			COALESCE(w.total_withdrawn, 0) AS total_withdrawn
		FROM 
			users u
		LEFT JOIN 
			accrued a ON u.id = a.user_id
		LEFT JOIN 
			withdrawn w ON u.id = w.user_id
		WHERE
			u.id = $1;
	`

	row := r.db.QueryRow(query, userID, OrderStatusProcessed)

	var dest struct {
		UserID         int     `db:"user_id"`
		CurrentBalance float64 `db:"current_balance"`
		TotalWithdrawn float64 `db:"total_withdrawn"`
	}

	err := row.Scan(&dest.UserID, &dest.CurrentBalance, &dest.TotalWithdrawn)

	if err != nil {
		return dtos.Balance{}, fmt.Errorf("%s: %w", op, err)
	}

	return dtos.Balance{UserID: dest.UserID, Current: dest.CurrentBalance, Withdrawn: dest.TotalWithdrawn}, nil
}

func (r *DBBalanceRepository) CreateWithdraw(ctx context.Context, userID int, orderID string, sum float64) (int, error) {
	op := "balanceRepo.createWithdraw"

	stmt := table.Withdraws.INSERT(table.Withdraws.Amount, table.Withdraws.UserID, table.Withdraws.OrderID).
		VALUES(sum, userID, orderID).
		RETURNING(table.Withdraws.ID)

	var dest model.Withdraws

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(dest.ID), nil
}

func (r *DBBalanceRepository) GetWithdrawalsByUser(ctx context.Context, userID int) ([]dtos.Withdraw, error) {
	op := "balanceRepo.getWithdrawalsByUser"

	stmt := table.Withdraws.SELECT(table.Withdraws.AllColumns).WHERE(table.Withdraws.UserID.EQ(postgres.Int64(int64(userID))))

	var dest []model.Withdraws

	err := stmt.QueryContext(ctx, r.db, &dest)

	result := make([]dtos.Withdraw, len(dest))

	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	for i, entity := range dest {
		result[i] = mapWithdrawnEntityToDto(entity)
	}

	return result, nil
}

func NewDBBalanceRepository(db *sql.DB) *DBBalanceRepository {
	return &DBBalanceRepository{db}
}

var _ BalanceRepository = &DBBalanceRepository{}

func mapWithdrawnEntityToDto(entity model.Withdraws) dtos.Withdraw {
	return dtos.Withdraw{
		ID:          int(entity.ID),
		UserID:      int(entity.UserID),
		OrderID:     entity.OrderID,
		Amount:      entity.Amount,
		ProcessedAt: entity.CreatedAt,
	}

}
