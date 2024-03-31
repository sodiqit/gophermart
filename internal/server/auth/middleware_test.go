package auth_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestJWTMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := chi.NewRouter()

	tokenServiceMock := auth.NewMockTokenService(ctrl)

	r.Use(auth.JWTAuth(tokenServiceMock))

	r.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		user := auth.ExtractUserFromContext(r.Context())

		w.Write([]byte(fmt.Sprintf("%x", user.ID)))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New().SetBaseURL(ts.URL)

	tests := []struct {
		name           string
		method         string
		url            string
		setupMock      func()
		header         string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "should return 401 if token not provide",
			method:         http.MethodGet,
			url:            "/test/",
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func() {},
			header:         "",
		},
		{
			name:           "should return 401 if token not provide in Bearer format",
			method:         http.MethodGet,
			url:            "/test/",
			expectedStatus: http.StatusUnauthorized,
			header:         "token",
			setupMock:      func() {},
		},
		{
			name:           "should return error if token not valid",
			method:         http.MethodGet,
			url:            "/test/",
			expectedStatus: http.StatusUnauthorized,
			header:         "Bearer token",
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate("token").Return(&auth.Claims{}, errors.New("invalid token"))
			},
		},
		{
			name:           "should return 200 if token valid",
			method:         http.MethodGet,
			url:            "/test/",
			expectedStatus: http.StatusOK,
			header:         "Bearer token",
			expectedResult: "2",
			setupMock: func() {
				tokenServiceMock.EXPECT().Validate(gomock.Any()).Return(&auth.Claims{TokenUser: auth.TokenUser{ID: 2}}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			req := client.R()

			req.Method = tc.method
			req.URL = tc.url

			if tc.header != "" {
				req.SetHeader("Authorization", tc.header)
			}

			resp, err := req.Send()

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, resp.StatusCode())
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedResult, resp.String())
			}
		})
	}
}
