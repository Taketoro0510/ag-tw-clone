package repository

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeRepository interface {
	Create(ctx context.Context, postID, userID string) error
	Delete(ctx context.Context, postID, userID string) error
	CountByPostID(ctx context.Context, postID string) (int, error)
	IsLikedByUser(ctx context.Context, postID, userID string) (bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, postID, userID string) error {
	m := &model.Like{
		PostID: postID,
		UserID: userID,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
}

func (r *likeRepository) Delete(ctx context.Context, postID, userID string) error {
	return r.db.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.Like{}).Error
}

func (r *likeRepository) CountByPostID(ctx context.Context, postID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}

func (r *likeRepository) IsLikedByUser(ctx context.Context, postID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Like{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	return count > 0, err
}
