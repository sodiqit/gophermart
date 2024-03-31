package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthController_handleRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := chi.NewRouter()

	authServiceMock := auth.NewMockAuthService(ctrl)
	logger := logger.New("info")

	c := auth.NewController(logger, authServiceMock)

	r.Mount("/", c.Route())

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
			url:            "/register/",
			body:           "test",
			expectedStatus: http.StatusUnsupportedMediaType,
			setupMock:      func() {},
		},
		{
			name:           "should success dto validate",
			method:         http.MethodPost,
			url:            "/register/",
			body:           `{"username": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
		{
			name:           "should success register new user",
			method:         http.MethodPost,
			url:            "/register/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				authServiceMock.EXPECT().Register(gomock.Any(), "test", "test").Return("test_token", nil)
			},
		},
		{
			name:           "should handle user already exist error",
			method:         http.MethodPost,
			url:            "/register/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusConflict,
			setupMock: func() {
				authServiceMock.EXPECT().Register(gomock.Any(), "test", "test").Return("", auth.ErrUserAlreadyExist)
			},
		},
		{
			name:           "should handle other auth service errors",
			method:         http.MethodPost,
			url:            "/register/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func() {
				authServiceMock.EXPECT().Register(gomock.Any(), "test", "test").Return("", errors.New("unexpected error"))
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

			resp, err := req.Send()

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, resp.StatusCode())
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, "Bearer test_token", resp.Header().Get("Authorization"))
			}
		})
	}
}

func TestAuthController_handleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := chi.NewRouter()

	authServiceMock := auth.NewMockAuthService(ctrl)
	logger := logger.New("info")

	c := auth.NewController(logger, authServiceMock)

	r.Mount("/", c.Route())

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
			url:            "/login/",
			body:           "test",
			expectedStatus: http.StatusUnsupportedMediaType,
			setupMock:      func() {},
		},
		{
			name:           "should success dto validate",
			method:         http.MethodPost,
			url:            "/login/",
			body:           `{"username": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
		{
			name:           "should success generate token for user",
			method:         http.MethodPost,
			url:            "/login/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				authServiceMock.EXPECT().Login(gomock.Any(), "test", "test").Return("test_token", nil)
			},
		},
		{
			name:           "should handle error if user not exist",
			method:         http.MethodPost,
			url:            "/login/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func() {
				authServiceMock.EXPECT().Login(gomock.Any(), "test", "test").Return("", auth.ErrUserNotFound)
			},
		},
		{
			name:           "should handle error if user pass not correct",
			method:         http.MethodPost,
			url:            "/login/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func() {
				authServiceMock.EXPECT().Login(gomock.Any(), "test", "test").Return("", auth.ErrIncorrectPassword)
			},
		},
		{
			name:           "should handle other auth service errors",
			method:         http.MethodPost,
			url:            "/login/",
			body:           `{"username": "test", "password": "test"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func() {
				authServiceMock.EXPECT().Login(gomock.Any(), "test", "test").Return("", errors.New("unexpected error"))
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

			resp, err := req.Send()

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, resp.StatusCode())
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, "Bearer test_token", resp.Header().Get("Authorization"))
			}
		})
	}
}
