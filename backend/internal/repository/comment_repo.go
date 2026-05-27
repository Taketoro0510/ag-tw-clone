package repository

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string, cursor string, limit int) ([]*domain.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
	GetCommentByID(ctx context.Context, commentID string) (*domain.Comment, error)
	LikeComment(ctx context.Context, commentID, userID string) error
	UnlikeComment(ctx context.Context, commentID, userID string) error
	GetLikeCount(ctx context.Context, commentID string) (int, error)
	IsLikedByMe(ctx context.Context, commentID, userID string) (bool, error)
	BookmarkComment(ctx context.Context, commentID, userID string) error
	UnbookmarkComment(ctx context.Context, commentID, userID string) error
	GetBookmarkCount(ctx context.Context, commentID string) (int, error)
	IsBookmarkedByMe(ctx context.Context, commentID, userID string) (bool, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	m := &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		AuthorID:  comment.AuthorID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string, cursor string, limit int) ([]*domain.Comment, error) {
	var models []model.Comment
	query := r.db.WithContext(ctx).Where("post_id = ? AND deleted_at IS NULL", postID).Preload("Author")

	if cursor != "" {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("id DESC").Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var comments []*domain.Comment
	for _, m := range models {
		comments = append(comments, r.toDomain(&m))
	}
	return comments, nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, commentID string) error {
	return r.db.WithContext(ctx).Where("id = ?", commentID).Delete(&model.Comment{}).Error
}

func (r *commentRepository) GetCommentByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	var m model.Comment
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", commentID).Preload("Author").First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPostNotFound // reusing not found
		}
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *commentRepository) LikeComment(ctx context.Context, commentID, userID string) error {
	m := &model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *commentRepository) UnlikeComment(ctx context.Context, commentID, userID string) error {
	return r.db.WithContext(ctx).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&model.CommentLike{}).Error
}

func (r *commentRepository) GetLikeCount(ctx context.Context, commentID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).Where("comment_id = ?", commentID).Count(&count).Error
	return int(count), err
}

func (r *commentRepository) IsLikedByMe(ctx context.Context, commentID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error
	return count > 0, err
}

func (r *commentRepository) BookmarkComment(ctx context.Context, commentID, userID string) error {
	m := &model.CommentBookmark{
		CommentID: commentID,
		UserID:    userID,
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *commentRepository) UnbookmarkComment(ctx context.Context, commentID, userID string) error {
	return r.db.WithContext(ctx).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&model.CommentBookmark{}).Error
}

func (r *commentRepository) GetBookmarkCount(ctx context.Context, commentID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentBookmark{}).Where("comment_id = ?", commentID).Count(&count).Error
	return int(count), err
}

func (r *commentRepository) IsBookmarkedByMe(ctx context.Context, commentID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentBookmark{}).Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error
	return count > 0, err
}

func (r *commentRepository) toDomain(m *model.Comment) *domain.Comment {
	dc := &domain.Comment{
		ID:        m.ID,
		PostID:    m.PostID,
		AuthorID:  m.AuthorID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
		DeletedAt: m.DeletedAt,
	}
	if m.Author.ID != "" {
		dc.Author = &domain.User{
			ID:          m.Author.ID,
			FirebaseUID: m.Author.FirebaseUID,
			Email:       m.Author.Email,
			DisplayName: m.Author.DisplayName,
			AvatarURL:   m.Author.AvatarURL,
			CreatedAt:   m.Author.CreatedAt,
			UpdatedAt:   m.Author.UpdatedAt,
		}
	}
	return dc
}
