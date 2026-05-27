package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockAuthUseCase struct {
	loginFn func(ctx context.Context, idToken string) (string, error)
}

func (m *mockAuthUseCase) Login(ctx context.Context, idToken string) (string, error) {
	return m.loginFn(ctx, idToken)
}

func TestAuthHandler_CreateSession(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		reqBody        string
		mockLogin      func(ctx context.Context, idToken string) (string, error)
		expectedStatus int
		expectedBody   string // partial string match
	}{
		{
			name:           "invalid body",
			reqBody:        `{invalid}`,
			mockLogin:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"リクエストボディが不正です"`,
		},
		{
			name:    "unauthorized token",
			reqBody: `{"idToken":"invalid-token"}`,
			mockLogin: func(ctx context.Context, idToken string) (string, error) {
				return "", domain.ErrUnauthorized
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `"message":"認証トークンが不正です"`,
		},
		{
			name:    "internal error",
			reqBody: `{"idToken":"valid-token"}`,
			mockLogin: func(ctx context.Context, idToken string) (string, error) {
				return "", errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"ログインに失敗しました"`,
		},
		{
			name:    "success",
			reqBody: `{"idToken":"valid-token"}`,
			mockLogin: func(ctx context.Context, idToken string) (string, error) {
				return "jwt-token-123", nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"token":"jwt-token-123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/sessions", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUC := &mockAuthUseCase{loginFn: tt.mockLogin}
			h := handler.NewAuthHandler(mockUC)

			err := h.CreateSession(c)

			// Echo handlers often return errors if unhandled, but here we return c.JSON
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}
