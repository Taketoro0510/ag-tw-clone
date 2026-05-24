package middleware

import (
	"net/http"
	"strings"

	"github.com/koitake1/cloudcode-sns/backend/internal/auth"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/labstack/echo/v4"
)

func Auth(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "missing or invalid token"},
				})
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := jwtService.VerifyToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "invalid or expired token"},
				})
			}
			c.Set("userID", userID)
			return next(c)
		}
	}
}
