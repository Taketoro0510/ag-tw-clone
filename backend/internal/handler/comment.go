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

type CommentHandler struct {
	commentUC usecase.CommentUseCase
}

func NewCommentHandler(cuc usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{commentUC: cuc}
}

func (h *CommentHandler) ListComments(c echo.Context) error {
	postID := c.Param("id")
	userID := c.Get("userID").(string)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	cursor := c.QueryParam("cursor")

	comments, err := h.commentUC.GetCommentsByPostID(c.Request().Context(), userID, postID, cursor, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "コメント一覧の取得に失敗しました"},
		})
	}

	items := make([]dto.Comment, len(comments))
	for i, c := range comments {
		items[i] = toCommentDTO(c)
	}
	var nextCursor *string
	if len(comments) > 0 {
		lastID := comments[len(comments)-1].ID
		nextCursor = &lastID
	}
	return c.JSON(http.StatusOK, dto.PaginatedCommentsResponse{Items: items, NextCursor: nextCursor})
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	postID := c.Param("id")
	userID := c.Get("userID").(string)
	var req dto.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "リクエストボディが不正です"},
		})
	}

	comment, err := h.commentUC.CreateComment(c.Request().Context(), userID, postID, req.Body)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "バリデーションに失敗しました"},
			})
		}
		if errors.Is(err, domain.ErrPostNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: "投稿が見つかりません"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "コメントの作成に失敗しました"},
		})
	}
	return c.JSON(http.StatusCreated, toCommentDTO(comment))
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	commentID := c.Param("commentId")
	userID := c.Get("userID").(string)

	err := h.commentUC.DeleteComment(c.Request().Context(), userID, commentID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "FORBIDDEN", Message: "他人のコメントは削除できません"},
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "コメントの削除に失敗しました"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CommentHandler) LikeComment(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.commentUC.LikeComment(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "いいねに失敗しました"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CommentHandler) UnlikeComment(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.commentUC.UnlikeComment(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "いいねの解除に失敗しました"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CommentHandler) BookmarkComment(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.commentUC.BookmarkComment(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "ブックマークに失敗しました"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CommentHandler) UnbookmarkComment(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(string)
	err := h.commentUC.UnbookmarkComment(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "ブックマーク解除に失敗しました"},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func toCommentDTO(c *domain.Comment) dto.Comment {
	var author *dto.User
	if c.Author != nil {
		u := toUserDTO(c.Author)
		author = &u
	}

	return dto.Comment{
		ID:             c.ID,
		PostID:         c.PostID,
		AuthorID:       c.AuthorID,
		Author:         author,
		Body:           c.Body,
		LikeCount:      c.LikeCount,
		LikedByMe:      c.LikedByMe,
		BookmarkCount:  c.BookmarkCount,
		BookmarkedByMe: c.BookmarkedByMe,
		CreatedAt:      c.CreatedAt,
	}
}
