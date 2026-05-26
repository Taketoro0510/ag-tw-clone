package handler

import (
	"net/http"
	"strconv"

	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type BookmarkHandler struct {
	bmUC usecase.BookmarkUseCase
}

func NewBookmarkHandler(uc usecase.BookmarkUseCase) *BookmarkHandler {
	return &BookmarkHandler{bmUC: uc}
}

// BookmarkPost godoc
// @Summary      Bookmark a post
// @Description  Bookmark a post
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /posts/{id}/bookmarks [post]
func (h *BookmarkHandler) BookmarkPost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.bmUC.BookmarkPost(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to bookmark post"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// UnbookmarkPost godoc
// @Summary      Unbookmark a post
// @Description  Unbookmark a post
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /posts/{id}/bookmarks [delete]
func (h *BookmarkHandler) UnbookmarkPost(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.bmUC.UnbookmarkPost(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to unbookmark post"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListBookmarks godoc
// @Summary      List bookmarked posts
// @Description  List bookmarked posts
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        limit  query     int     false "Limit" default(20)
// @Param        cursor query     string  false "Cursor"
// @Security     BearerAuth
// @Success      200  {object}  dto.PaginatedPostsResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /bookmarks [get]
func (h *BookmarkHandler) ListBookmarks(c echo.Context) error {
	userID := c.Get("userID").(string)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	cursor := c.QueryParam("cursor")

	posts, err := h.bmUC.ListBookmarks(c.Request().Context(), userID, cursor, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to list bookmarked posts"},
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
