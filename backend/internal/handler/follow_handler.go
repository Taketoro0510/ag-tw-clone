package handler

import (
	"errors"
	"net/http"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type FollowHandler struct {
	followUC usecase.FollowUseCase
}

func NewFollowHandler(uc usecase.FollowUseCase) *FollowHandler {
	return &FollowHandler{followUC: uc}
}

// FollowUser godoc
// @Summary      Follow a user
// @Description  Follow a user by ID
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID to follow"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /users/{id}/follow [post]
func (h *FollowHandler) FollowUser(c echo.Context) error {
	followeeID := c.Param("id")
	followerID := c.Get("userID").(string)

	err := h.followUC.FollowUser(c.Request().Context(), followerID, followeeID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: "user not found"},
			})
		}
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: err.Error()},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// UnfollowUser godoc
// @Summary      Unfollow a user
// @Description  Unfollow a user by ID
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID to unfollow"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /users/{id}/follow [delete]
func (h *FollowHandler) UnfollowUser(c echo.Context) error {
	followeeID := c.Param("id")
	followerID := c.Get("userID").(string)

	err := h.followUC.UnfollowUser(c.Request().Context(), followerID, followeeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to unfollow user"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListFollowers godoc
// @Summary      List followers
// @Description  Get users who follow this user
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {array}   dto.User
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /users/{id}/followers [get]
func (h *FollowHandler) ListFollowers(c echo.Context) error {
	userID := c.Param("id")
	users, err := h.followUC.GetFollowers(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to get followers"},
		})
	}

	dtos := make([]dto.User, len(users))
	for i, u := range users {
		dtos[i] = toUserDTO(u)
	}
	return c.JSON(http.StatusOK, dtos)
}

// ListFollowings godoc
// @Summary      List followings
// @Description  Get users followed by this user
// @Tags         follows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {array}   dto.User
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /users/{id}/following [get]
func (h *FollowHandler) ListFollowings(c echo.Context) error {
	userID := c.Param("id")
	users, err := h.followUC.GetFollowings(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to get followed users"},
		})
	}

	dtos := make([]dto.User, len(users))
	for i, u := range users {
		dtos[i] = toUserDTO(u)
	}
	return c.JSON(http.StatusOK, dtos)
}
