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

type PostHandler struct {
	postUC usecase.PostUseCase
	likeUC usecase.LikeUseCase
}

func NewPostHandler(puc usecase.PostUseCase, luc usecase.LikeUseCase) *PostHandler {
	return &PostHandler{postUC: puc, likeUC: luc}
}

// ListPosts godoc
// @Summary      List posts
// @Description  Get global timeline
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        limit  query     int     false "Limit" default(20)
// @Param        cursor query     string  false "Cursor"
// @Security     BearerAuth
// @Success      200  {object}  dto.PaginatedPostsResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /posts [get]
func (h *PostHandler) ListPosts(c echo.Context) error {
	userID := c.Get("userID").(string)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	cursor := c.QueryParam("cursor")

	posts, err := h.postUC.GetTimeline(c.Request().Context(), userID, cursor, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to list posts"},
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

// CreatePost godoc
// @Summary      Create a post
// @Description  Create a new post with text and optional media
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePostRequest true "Post info"
// @Security     BearerAuth
// @Success      201  {object}  dto.Post
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      429  {object}  dto.ErrorResponse
// @Router       /posts [post]
func (h *PostHandler) CreatePost(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req dto.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "invalid body"},
		})
	}

	post, err := h.postUC.CreatePost(c.Request().Context(), userID, req.Body, req.MediaType, req.MediaPath, nil)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "validation failed"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to create post"},
		})
	}
	return c.JSON(http.StatusCreated, toPostDTO(post))
}

// GetPost godoc
// @Summary      Get a post
// @Description  Get post details
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      200  {object}  dto.Post
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /posts/{id} [get]
func (h *PostHandler) GetPost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	post, err := h.postUC.GetPost(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "RESOURCE_NOT_FOUND", Message: "post not found"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to get post"},
		})
	}
	return c.JSON(http.StatusOK, toPostDTO(post))
}

// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete own post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /posts/{id} [delete]
func (h *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.postUC.DeletePost(c.Request().Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "FORBIDDEN", Message: "cannot delete others post"},
			})
		}
		if errors.Is(err, domain.ErrPostNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "RESOURCE_NOT_FOUND", Message: "post not found"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to delete post"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// LikePost godoc
// @Summary      Like a post
// @Description  Like a post
// @Tags         likes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      429  {object}  dto.ErrorResponse
// @Router       /posts/{id}/likes [post]
func (h *PostHandler) LikePost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.likeUC.LikePost(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to like post"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// UnlikePost godoc
// @Summary      Unlike a post
// @Description  Unlike a post
// @Tags         likes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      429  {object}  dto.ErrorResponse
// @Router       /posts/{id}/likes [delete]
func (h *PostHandler) UnlikePost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.likeUC.UnlikePost(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to unlike post"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}
