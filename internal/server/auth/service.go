package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username string, password string) (string, error)
	Login(ctx context.Context, username string, password string) (string, error)
}

var ErrUserAlreadyExist = errors.New("user already exist")
var ErrUserNotFound = errors.New("user not found")
var ErrIncorrectPassword = errors.New("incorrect password")

type SimpleAuthService struct {
	tokenService TokenService
	userRepo     repository.UserRepository
}

func (s *SimpleAuthService) Register(ctx context.Context, username string, password string) (string, error) {
	op := "authService.register"

	exist, err := s.userRepo.Exist(ctx, username)

	if err != nil {
		return "", err
	}

	if exist {
		return "", fmt.Errorf("%s: %w", op, ErrUserAlreadyExist)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	userId, err := s.userRepo.Create(ctx, dtos.User{Login: username, PasswordHash: string(passHash)})

	if err != nil {
		return "", err
	}

	return s.tokenService.Build(userId)
}

func (s *SimpleAuthService) Login(ctx context.Context, username string, password string) (string, error) {
	op := "authService.login"

	user, err := s.userRepo.FindByLogin(ctx, username)

	if errors.Is(err, qrm.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrIncorrectPassword)
	}

	return s.tokenService.Build(user.ID)
}

func NewSimpleAuthService(tokenService TokenService, userRepo repository.UserRepository) *SimpleAuthService {
	return &SimpleAuthService{
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}
