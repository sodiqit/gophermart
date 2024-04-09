package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/sodiqit/gophermart/gen/gophermart_db/public/model"
	"github.com/sodiqit/gophermart/gen/gophermart_db/public/table"
	"github.com/sodiqit/gophermart/internal/server/dtos"
)

type UserRepository interface {
	Create(ctx context.Context, user dtos.User) (int, error)
	FindByLogin(ctx context.Context, username string) (dtos.User, error)
	Exist(ctx context.Context, login string) (bool, error)
}

type DBUserRepository struct {
	db *sql.DB
}

func (r *DBUserRepository) Create(ctx context.Context, user dtos.User) (int, error) {
	op := "userRepo.create"

	stmt := table.Users.INSERT(table.Users.Login, table.Users.PasswordHash).VALUES(user.Login, user.PasswordHash).RETURNING(table.Users.ID)

	var dest model.Users

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(dest.ID), err
}

func (r *DBUserRepository) FindByLogin(ctx context.Context, username string) (dtos.User, error) {
	op := "userRepo.findOne"

	stmt := table.Users.SELECT(table.Users.ID, table.Users.Login, table.Users.PasswordHash, table.Users.CreatedAt).WHERE(table.Users.Login.EQ(postgres.String(username)))

	var dest model.Users

	err := stmt.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return dtos.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return mapUserEntityToDto(dest), nil
}

func (r *DBUserRepository) Exist(ctx context.Context, login string) (bool, error) {
	op := "userRepo.exist"

	stmt := table.Users.SELECT(table.Users.ID).WHERE(table.Users.Login.EQ(postgres.String(login)))

	existStmt := postgres.SELECT(postgres.EXISTS(stmt)).FROM(table.Users)

	sqlStmt, args := existStmt.Sql()

	var exist bool

	err := r.db.QueryRow(sqlStmt, args...).Scan(&exist)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exist, nil
}

func mapUserEntityToDto(userEntity model.Users) dtos.User {
	return dtos.User{
		ID:           int(userEntity.ID),
		Login:        userEntity.Login,
		PasswordHash: userEntity.PasswordHash,
		CreatedAt:    userEntity.CreatedAt,
	}

}

var _ UserRepository = (*DBUserRepository)(nil)

func NewDBUserRepository(db *sql.DB) *DBUserRepository {
	return &DBUserRepository{db: db}
}
