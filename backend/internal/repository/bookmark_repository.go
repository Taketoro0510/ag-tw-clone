package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	Bookmark(ctx context.Context, userID, postID string) error
	Unbookmark(ctx context.Context, userID, postID string) error
	CountByPostID(ctx context.Context, postID string) (int, error)
	IsBookmarkedByUser(ctx context.Context, postID, userID string) (bool, error)
	ListBookmarkedPosts(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error)
}

type bookmarkRepository struct {
	db       *gorm.DB
	postRepo PostRepository
}

func NewBookmarkRepository(db *gorm.DB, postRepo PostRepository) BookmarkRepository {
	return &bookmarkRepository{db: db, postRepo: postRepo}
}

func (r *bookmarkRepository) Bookmark(ctx context.Context, userID, postID string) error {
	id, _ := uuid.NewV7()
	bm := model.Bookmark{
		ID:     id.String(),
		UserID: userID,
		PostID: postID,
	}
	return r.db.WithContext(ctx).Create(&bm).Error
}

func (r *bookmarkRepository) Unbookmark(ctx context.Context, userID, postID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Bookmark{}).Error
}

func (r *bookmarkRepository) CountByPostID(ctx context.Context, postID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}

func (r *bookmarkRepository) IsBookmarkedByUser(ctx context.Context, postID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	return count > 0, err
}

func (r *bookmarkRepository) ListBookmarkedPosts(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error) {
	var ms []model.Post
	q := r.db.WithContext(ctx).
		Preload("Author").
		Joins("JOIN bookmarks ON bookmarks.post_id = posts.id").
		Where("bookmarks.user_id = ? AND posts.deleted_at IS NULL", userID)

	if cursor != "" {
		q = q.Where("posts.id < ?", cursor)
	}

	err := q.Order("posts.id DESC").Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}

	// Because of circular deps or duplicating toDomain, let's just implement a quick mapper
	// or we can use a helper method. For MVP, we can map directly here.
	res := make([]*domain.Post, len(ms))
	for i, m := range ms {
		p := &domain.Post{
			ID:        m.ID,
			AuthorID:  m.AuthorID,
			Body:      m.Body,
			MediaType: m.MediaType,
			MediaPath: m.MediaPath,
			MediaURL:  m.MediaURL,
			CreatedAt: m.CreatedAt,
			DeletedAt: m.DeletedAt,
		}
		if m.Author.ID != "" {
			p.Author = &domain.User{
				ID:          m.Author.ID,
				FirebaseUID: m.Author.FirebaseUID,
				Email:       m.Author.Email,
				DisplayName: m.Author.DisplayName,
				AvatarURL:   m.Author.AvatarURL,
				CreatedAt:   m.Author.CreatedAt,
				UpdatedAt:   m.Author.UpdatedAt,
			}
		}
		res[i] = p
	}
	return res, nil
}
