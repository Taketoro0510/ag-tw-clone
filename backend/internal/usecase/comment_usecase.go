package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, userID, postID, body string) (*domain.Comment, error)
	GetCommentsByPostID(ctx context.Context, userID, postID string, cursor string, limit int) ([]*domain.Comment, error)
	DeleteComment(ctx context.Context, userID, commentID string) error
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error
	BookmarkComment(ctx context.Context, userID, commentID string) error
	UnbookmarkComment(ctx context.Context, userID, commentID string) error
}

type commentUseCase struct {
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
}

func NewCommentUseCase(cr repository.CommentRepository, ur repository.UserRepository, pr repository.PostRepository) CommentUseCase {
	return &commentUseCase{commentRepo: cr, userRepo: ur, postRepo: pr}
}

func (u *commentUseCase) CreateComment(ctx context.Context, userID, postID, body string) (*domain.Comment, error) {
	if body == "" {
		return nil, domain.ErrValidation
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, domain.ErrPostNotFound
	}

	id, _ := uuid.NewV7()
	comment := &domain.Comment{
		ID:        id.String(),
		PostID:    postID,
		AuthorID:  user.ID,
		Author:    user,
		Body:      body,
		CreatedAt: time.Now(),
	}

	err = u.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Add counts
	comment.LikeCount = 0
	comment.LikedByMe = false
	comment.BookmarkCount = 0
	comment.BookmarkedByMe = false

	return comment, nil
}

func (u *commentUseCase) GetCommentsByPostID(ctx context.Context, userID, postID string, cursor string, limit int) ([]*domain.Comment, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	comments, err := u.commentRepo.GetCommentsByPostID(ctx, postID, cursor, limit)
	if err != nil {
		return nil, err
	}

	for _, c := range comments {
		c.LikeCount, _ = u.commentRepo.GetLikeCount(ctx, c.ID)
		c.LikedByMe, _ = u.commentRepo.IsLikedByMe(ctx, c.ID, user.ID)
		c.BookmarkCount, _ = u.commentRepo.GetBookmarkCount(ctx, c.ID)
		c.BookmarkedByMe, _ = u.commentRepo.IsBookmarkedByMe(ctx, c.ID, user.ID)
	}

	return comments, nil
}

func (u *commentUseCase) DeleteComment(ctx context.Context, userID, commentID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	comment, err := u.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.AuthorID != user.ID {
		return domain.ErrForbidden
	}

	return u.commentRepo.DeleteComment(ctx, commentID)
}

func (u *commentUseCase) LikeComment(ctx context.Context, userID, commentID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	return u.commentRepo.LikeComment(ctx, commentID, user.ID)
}

func (u *commentUseCase) UnlikeComment(ctx context.Context, userID, commentID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	return u.commentRepo.UnlikeComment(ctx, commentID, user.ID)
}

func (u *commentUseCase) BookmarkComment(ctx context.Context, userID, commentID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	return u.commentRepo.BookmarkComment(ctx, commentID, user.ID)
}

func (u *commentUseCase) UnbookmarkComment(ctx context.Context, userID, commentID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	return u.commentRepo.UnbookmarkComment(ctx, commentID, user.ID)
}
