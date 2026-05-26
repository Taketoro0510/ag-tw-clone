package repository

import (
	"context"
	"errors"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id string) (*domain.Post, error)
	Delete(ctx context.Context, id string) error
	ListGlobal(ctx context.Context, cursor string, limit int) ([]*domain.Post, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) toDomain(m *model.Post) *domain.Post {
	dp := &domain.Post{
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
		dp.Author = &domain.User{
			ID:          m.Author.ID,
			FirebaseUID: m.Author.FirebaseUID,
			Email:       m.Author.Email,
			DisplayName: m.Author.DisplayName,
			AvatarURL:   m.Author.AvatarURL,
			CreatedAt:   m.Author.CreatedAt,
			UpdatedAt:   m.Author.UpdatedAt,
		}
	}
	return dp
}

func (r *postRepository) toModel(d *domain.Post) *model.Post {
	return &model.Post{
		ID:        d.ID,
		AuthorID:  d.AuthorID,
		Body:      d.Body,
		MediaType: d.MediaType,
		MediaPath: d.MediaPath,
		MediaURL:  d.MediaURL,
		CreatedAt: d.CreatedAt,
		DeletedAt: d.DeletedAt,
	}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	m := r.toModel(post)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *postRepository) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	var m model.Post
	if err := r.db.WithContext(ctx).Preload("Author").Where("id = ? AND deleted_at IS NULL", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ListGlobal uses cursor based on (created_at, id) since UUID v7 has time-based ordering
func (r *postRepository) ListGlobal(ctx context.Context, cursor string, limit int) ([]*domain.Post, error) {
	var ms []model.Post
	q := r.db.WithContext(ctx).Preload("Author").Where("deleted_at IS NULL")
	if cursor != "" {
		// MVP: Cursor is just the ID because UUID v7 is naturally sortable
		q = q.Where("id < ?", cursor)
	}
	err := q.Order("id DESC").Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	res := make([]*domain.Post, len(ms))
	for i, m := range ms {
		res[i] = r.toDomain(&m)
	}
	return res, nil
}

func (r *postRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error) {
	var ms []model.Post
	q := r.db.WithContext(ctx).Preload("Author").Where("author_id = ? AND deleted_at IS NULL", userID)
	if cursor != "" {
		q = q.Where("id < ?", cursor)
	}
	err := q.Order("id DESC").Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	res := make([]*domain.Post, len(ms))
	for i, m := range ms {
		res[i] = r.toDomain(&m)
	}
	return res, nil
}
