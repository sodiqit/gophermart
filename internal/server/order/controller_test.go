package order_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/dtos"
	"github.com/sodiqit/gophermart/internal/server/order"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOrderController_handleOrderUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := chi.NewRouter()

	orderServiceMock := order.NewMockOrderService(ctrl)
	tokenServiceMock := auth.NewMockTokenService(ctrl)
	logger := logger.New("info")

	c := order.NewController(logger, tokenServiceMock, orderServiceMock)

	r.Mount("/orders", c.Route())

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New().SetBaseURL(ts.URL)

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		contentType    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:           "invalid Content-type",
			method:         http.MethodPost,
			url:            "/orders",
			body:           `{"test": true}`,
			contentType:    "application/json",
			expectedStatus: http.StatusUnsupportedMediaType,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{}, nil)
			},
		},
		{
			name:           "should return 401 if token invalid",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678904",
			contentType:    "text/plain",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{}, errors.New("invalid"))
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:           "should success validate order number by Luhn algorithm",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678904",
			contentType:    "text/plain",
			expectedStatus: http.StatusUnprocessableEntity,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:           "should success upload new order",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678903",
			contentType:    "text/plain",
			expectedStatus: http.StatusCreated,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name:           "should return 200 if user already add order",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678903",
			contentType:    "text/plain",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(order.ErrUserAlreadyUploadOrder)
			},
		},
		{
			name:           "should return 409 if another user already add order",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678903",
			contentType:    "text/plain",
			expectedStatus: http.StatusConflict,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(order.ErrOrderAlreadyUploadByAnotherUser)
			},
		},
		{
			name:           "should handle unexpected error",
			method:         http.MethodPost,
			url:            "/orders",
			body:           "12345678903",
			contentType:    "text/plain",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orderServiceMock.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("unexpected error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			req := client.R().SetBody(tc.body)

			req.Method = tc.method
			req.URL = tc.url

			if tc.contentType != "" {
				req.SetHeader("Content-Type", tc.contentType)
			}

			req.SetHeader("Authorization", "Bearer test")

			resp, err := req.Send()

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, resp.StatusCode())
		})
	}
}

func TestOrderController_handleGetUserList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := chi.NewRouter()

	orderServiceMock := order.NewMockOrderService(ctrl)
	tokenServiceMock := auth.NewMockTokenService(ctrl)
	logger := logger.New("info")

	c := order.NewController(logger, tokenServiceMock, orderServiceMock)

	r.Mount("/orders", c.Route())

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New().SetBaseURL(ts.URL)

	now := time.Now().UTC()

	tests := []struct {
		name           string
		method         string
		url            string
		contentType    string
		setupMock      func()
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "should return 401 if token invalid",
			method:         http.MethodGet,
			url:            "/orders",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{}, errors.New("invalid"))
				orderServiceMock.EXPECT().GetUserOrders(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:           "should success return user orders",
			method:         http.MethodGet,
			url:            "/orders",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedResult: `[{"number": "1234", "status": "NEW", "uploaded_at": "` + now.Format(time.RFC3339Nano) + `"}]`,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orders := make([]dtos.Order, 1)
				orders[0] = dtos.Order{ID: "1234", Status: repository.OrderStatusNew, CreatedAt: now}
				orderServiceMock.EXPECT().GetUserOrders(gomock.Any(), 1).Return(orders, nil)
			},
		},
		{
			name:           "should return 204 if no content",
			method:         http.MethodGet,
			url:            "/orders",
			expectedStatus: http.StatusNoContent,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orders := make([]dtos.Order, 0)
				orderServiceMock.EXPECT().GetUserOrders(gomock.Any(), 1).Return(orders, nil)
			},
		},
		{
			name:           "should return 204 if no content",
			method:         http.MethodGet,
			url:            "/orders",
			expectedStatus: http.StatusNoContent,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orders := make([]dtos.Order, 0)
				orderServiceMock.EXPECT().GetUserOrders(gomock.Any(), 1).Return(orders, nil)
			},
		},
		{
			name:           "should handle unexpected error",
			method:         http.MethodGet,
			url:            "/orders",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 1}}, nil)
				orders := make([]dtos.Order, 0)
				orderServiceMock.EXPECT().GetUserOrders(gomock.Any(), 1).Return(orders, errors.New("error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			req := client.R()

			req.Method = tc.method
			req.URL = tc.url

			req.SetHeader("Authorization", "Bearer test")

			resp, err := req.Send()

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, resp.StatusCode())
			if tc.expectedStatus == http.StatusOK {
				require.JSONEq(t, tc.expectedResult, resp.String())
			}
		})
	}
}
