package repository

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followeeID string) error
	Unfollow(ctx context.Context, followerID, followeeID string) error
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	CountFollowers(ctx context.Context, userID string) (int, error)
	CountFollowings(ctx context.Context, userID string) (int, error)
	ListFollowers(ctx context.Context, userID string) ([]*domain.User, error)
	ListFollowings(ctx context.Context, userID string) ([]*domain.User, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	m := &model.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).Error
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepository) CountFollowers(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("followee_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

func (r *followRepository) CountFollowings(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

func (r *followRepository) ListFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	var ms []model.User
	err := r.db.WithContext(ctx).
		Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.followee_id = ?", userID).
		Order("follows.created_at DESC").
		Find(&ms).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.User, len(ms))
	for i, m := range ms {
		res[i] = &domain.User{
			ID:          m.ID,
			FirebaseUID: m.FirebaseUID,
			Email:       m.Email,
			DisplayName: m.DisplayName,
			AvatarURL:   m.AvatarURL,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return res, nil
}

func (r *followRepository) ListFollowings(ctx context.Context, userID string) ([]*domain.User, error) {
	var ms []model.User
	err := r.db.WithContext(ctx).
		Joins("JOIN follows ON follows.followee_id = users.id").
		Where("follows.follower_id = ?", userID).
		Order("follows.created_at DESC").
		Find(&ms).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.User, len(ms))
	for i, m := range ms {
		res[i] = &domain.User{
			ID:          m.ID,
			FirebaseUID: m.FirebaseUID,
			Email:       m.Email,
			DisplayName: m.DisplayName,
			AvatarURL:   m.AvatarURL,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return res, nil
}
