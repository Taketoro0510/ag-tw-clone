package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koitake1/cloudcode-sns/backend/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockJWTService struct {
	verifyTokenFn func(token string) (string, error)
}

func (m *mockJWTService) GenerateToken(userID string) (string, error) {
	return "", nil
}

func (m *mockJWTService) VerifyToken(tokenString string) (string, error) {
	return m.verifyTokenFn(tokenString)
}

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		authHeader     string
		mockVerify     func(token string) (string, error)
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "Header missing",
			authHeader:     "",
			mockVerify:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid header prefix",
			authHeader:     "Token something",
			mockVerify:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid token",
			authHeader: "Bearer invalidtoken",
			mockVerify: func(token string) (string, error) {
				return "", errors.New("invalid token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Valid token",
			authHeader: "Bearer validtoken",
			mockVerify: func(token string) (string, error) {
				return "user-123", nil
			},
			expectedStatus: http.StatusOK,
			expectedUserID: "user-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuth := &mockJWTService{verifyTokenFn: tt.mockVerify}
			middlewareFunc := middleware.Auth(mockAuth)

			handler := middlewareFunc(func(c echo.Context) error {
				userID := c.Get("userID").(string)
				assert.Equal(t, tt.expectedUserID, userID)
				return c.String(http.StatusOK, "ok")
			})

			_ = handler(c)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
