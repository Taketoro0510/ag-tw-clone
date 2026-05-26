package usecase

import (
	"context"
	"errors"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
)

type FollowUseCase interface {
	FollowUser(ctx context.Context, followerID, followeeID string) error
	UnfollowUser(ctx context.Context, followerID, followeeID string) error
	GetFollowCounts(ctx context.Context, userID string) (followers int, following int, err error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	GetFollowers(ctx context.Context, userID string) ([]*domain.User, error)
	GetFollowings(ctx context.Context, userID string) ([]*domain.User, error)
}

type followUseCase struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

func NewFollowUseCase(fr repository.FollowRepository, ur repository.UserRepository) FollowUseCase {
	return &followUseCase{
		followRepo: fr,
		userRepo:   ur,
	}
}

func (u *followUseCase) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	// Verify followee exists
	_, err := u.userRepo.FindByID(ctx, followeeID)
	if err != nil {
		return domain.ErrUserNotFound
	}
	return u.followRepo.Follow(ctx, followerID, followeeID)
}

func (u *followUseCase) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return u.followRepo.Unfollow(ctx, followerID, followeeID)
}

func (u *followUseCase) GetFollowCounts(ctx context.Context, userID string) (followers int, following int, err error) {
	followers, err = u.followRepo.CountFollowers(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	following, err = u.followRepo.CountFollowings(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	return followers, following, nil
}

func (u *followUseCase) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return u.followRepo.IsFollowing(ctx, followerID, followeeID)
}

func (u *followUseCase) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	return u.followRepo.ListFollowers(ctx, userID)
}

func (u *followUseCase) GetFollowings(ctx context.Context, userID string) ([]*domain.User, error) {
	return u.followRepo.ListFollowings(ctx, userID)
}
