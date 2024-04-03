package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/order"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestOrderService_upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepoMock := repository.NewMockOrderRepository(ctrl)

	s := order.NewSimpleOrderService(orderRepoMock)

	tests := []struct {
		name           string
		setupMock      func()
		userID         int
		expectedResult string
		expectedError  error
		wantErr        bool
	}{
		{
			name: "should return error if user already upload order",
			setupMock: func() {
				orderRepoMock.EXPECT().FindByOrderNumber(gomock.Any(), gomock.Any()).Return(dtos.Order{UserID: 1}, nil)
				orderRepoMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:        1,
			wantErr:       true,
			expectedError: order.ErrUserAlreadyUploadOrder,
		},
		{
			name: "should return error if another user upload order",
			setupMock: func() {
				orderRepoMock.EXPECT().FindByOrderNumber(gomock.Any(), gomock.Any()).Return(dtos.Order{UserID: 1}, nil)
				orderRepoMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			userID:        2,
			wantErr:       true,
			expectedError: order.ErrOrderAlreadyUploadByAnotherUser,
		},
		{
			name: "should return error if find by order order not correct execute",
			setupMock: func() {
				orderRepoMock.EXPECT().FindByOrderNumber(gomock.Any(), gomock.Any()).Return(dtos.Order{}, errors.New("unexpected error"))
				orderRepoMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: true,
		},
		{
			name: "should success upload new order",
			setupMock: func() {
				orderRepoMock.EXPECT().FindByOrderNumber(gomock.Any(), gomock.Any()).Return(dtos.Order{}, repository.ErrOrderNotFound)
				orderRepoMock.EXPECT().Create(gomock.Any(), 0, "1234", repository.OrderStatusNew).Times(1).Return("1234", nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := s.Upload(context.Background(), tc.userID, "1234")

			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			}

			if tc.wantErr {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
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
