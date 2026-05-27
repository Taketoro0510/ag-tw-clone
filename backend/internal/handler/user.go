package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func toUserDTO(u *domain.User) dto.User {
	return dto.User{
		ID:             u.ID,
		FirebaseUID:    u.FirebaseUID,
		Email:          u.Email,
		DisplayName:    u.DisplayName,
		AvatarURL:      u.AvatarURL,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		FollowedByMe:   u.FollowedByMe,
	}
}

// GetMe godoc
// @Summary      Get current user
// @Description  Returns the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.User
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("userID").(string)
	user, err := h.userUC.GetProfile(c.Request().Context(), userID, userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "user not found"},
		})
	}
	return c.JSON(http.StatusOK, toUserDTO(user))
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Returns user basic info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  dto.User
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	user, err := h.userUC.GetProfile(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "RESOURCE_NOT_FOUND", Message: "user not found"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal error"},
		})
	}
	return c.JSON(http.StatusOK, toUserDTO(user))
}

// ListUserPosts godoc
// @Summary      List user posts
// @Description  Get posts by a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "User ID"
// @Param        limit  query     int     false "Limit" default(20)
// @Param        cursor query     string  false "Cursor"
// @Security     BearerAuth
// @Success      200  {object}  dto.PaginatedPostsResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /users/{id}/posts [get]
func (h *UserHandler) ListUserPosts(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	cursor := c.QueryParam("cursor")

	posts, err := h.userUC.GetUserPosts(c.Request().Context(), id, userID, cursor, limit)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "RESOURCE_NOT_FOUND", Message: "user not found"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal error"},
		})
	}

	items := make([]dto.Post, len(posts))
	for i, p := range posts {
		items[i] = toPostDTO(p)
	}
	var nextCursor *string
	if len(posts) > 0 {
		lastID := posts[len(posts)-1].ID
		nextCursor = &lastID
	}
	return c.JSON(http.StatusOK, dto.PaginatedPostsResponse{Items: items, NextCursor: nextCursor})
}

func toPostDTO(p *domain.Post) dto.Post {
	var author *dto.User
	if p.Author != nil {
		u := toUserDTO(p.Author)
		author = &u
	}
	return dto.Post{
		ID:             p.ID,
		AuthorID:       p.AuthorID,
		Author:         author,
		Body:           p.Body,
		MediaType:      p.MediaType,
		MediaPath:      p.MediaPath,
		MediaURL:       p.MediaURL,
		LikeCount:      p.LikeCount,
		LikedByMe:      p.LikedByMe,
		BookmarkCount:  p.BookmarkCount,
		BookmarkedByMe: p.BookmarkedByMe,
		CreatedAt:      p.CreatedAt,
	}
}
