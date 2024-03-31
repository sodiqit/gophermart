package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenServiceMock := auth.NewMockTokenService(ctrl)
	userRepoMock := repository.NewMockUserRepository(ctrl)

	s := auth.NewSimpleAuthService(tokenServiceMock, userRepoMock)

	tests := []struct {
		name           string
		setupMock      func()
		expectedResult string
		expectedError  error
		wantErr        bool
	}{
		{
			name: "should return error if user already register",
			setupMock: func() {
				userRepoMock.EXPECT().Exist(gomock.Any(), "test").Return(true, nil)
			},
			wantErr:       true,
			expectedError: auth.ErrUserAlreadyExist,
		},
		{
			name: "should return error if create failed",
			setupMock: func() {
				userRepoMock.EXPECT().Exist(gomock.Any(), "test").Return(false, nil)
				userRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(0, errors.New("create failed"))
				tokenServiceMock.EXPECT().Build(gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "should success create user and generate token",
			setupMock: func() {
				userRepoMock.EXPECT().Exist(gomock.Any(), "test").Return(false, nil)
				userRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)
				tokenServiceMock.EXPECT().Build(1).Return("test_token", nil)
			},
			wantErr:        false,
			expectedResult: "test_token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			token, err := s.Register(context.Background(), "test", "test")

			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			}

			if tc.wantErr {
				require.NotNil(t, err)
			}

			if !tc.wantErr {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, token)
			}
		})
	}
}

func TestAuthService_login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenServiceMock := auth.NewMockTokenService(ctrl)
	userRepoMock := repository.NewMockUserRepository(ctrl)

	s := auth.NewSimpleAuthService(tokenServiceMock, userRepoMock)

	tests := []struct {
		name           string
		setupMock      func()
		expectedResult string
		expectedError  error
		wantErr        bool
	}{
		{
			name: "should return error if user not found",
			setupMock: func() {
				userRepoMock.EXPECT().FindByLogin(gomock.Any(), "test").Return(dtos.User{}, qrm.ErrNoRows)
			},
			wantErr:       true,
			expectedError: auth.ErrUserNotFound,
		},
		{
			name: "should return error if password not correct",
			setupMock: func() {
				passHash, _ := bcrypt.GenerateFromPassword([]byte("incorrect"), bcrypt.DefaultCost)
				userRepoMock.EXPECT().FindByLogin(gomock.Any(), "test").Return(dtos.User{ID: 1, Login: "test", PasswordHash: string(passHash)}, nil)
				tokenServiceMock.EXPECT().Build(gomock.Any()).Times(0)
			},
			wantErr:       true,
			expectedError: auth.ErrIncorrectPassword,
		},
		{
			name: "should success generate token",
			setupMock: func() {
				passHash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
				userRepoMock.EXPECT().FindByLogin(gomock.Any(), "test").Return(dtos.User{ID: 1, Login: "test", PasswordHash: string(passHash)}, nil)
				tokenServiceMock.EXPECT().Build(1).Return("test_token", nil)
			},
			wantErr:        false,
			expectedResult: "test_token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			token, err := s.Login(context.Background(), "test", "test")

			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			}

			if tc.wantErr {
				require.NotNil(t, err)
			}

			if !tc.wantErr {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, token)
			}
		})
	}
}
