package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username string, password string) (string, error)
	LogIn(ctx context.Context, username string, password string) (string, error)
}

var ErrUserAlreadyExist = errors.New("user already exist")

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

func (s *SimpleAuthService) LogIn(ctx context.Context, username string, password string) (string, error) {
	return "", nil
}

func NewSimpleAuthService(tokenService TokenService, userRepo repository.UserRepository) *SimpleAuthService {
	return &SimpleAuthService{
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}
