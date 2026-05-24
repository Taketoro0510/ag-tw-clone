package handler

import (
	"errors"
	"net/http"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: uc}
}

// CreateSession godoc
// @Summary      Create a new session
// @Description  Exchanges Firebase ID token for a custom JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSessionRequest true "ID Token"
// @Success      200  {object}  dto.CreateSessionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /auth/sessions [post]
func (h *AuthHandler) CreateSession(c echo.Context) error {
	var req dto.CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "invalid body"},
		})
	}
	token, err := h.authUC.Login(c.Request().Context(), req.IdToken)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "invalid id token"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to login"},
		})
	}
	return c.JSON(http.StatusOK, dto.CreateSessionResponse{Token: token})
}
