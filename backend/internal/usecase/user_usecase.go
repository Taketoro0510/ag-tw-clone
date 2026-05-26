package usecase

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
)

type UserUseCase interface {
	GetProfile(ctx context.Context, userID, requestingUserID string) (*domain.User, error)
	GetUserPosts(ctx context.Context, userID, requestingUserID string, cursor string, limit int) ([]*domain.Post, error)
}

type userUseCase struct {
	userRepo   repository.UserRepository
	postRepo   repository.PostRepository
	likeRepo   repository.LikeRepository
	bmRepo     repository.BookmarkRepository
	followRepo repository.FollowRepository
}

func NewUserUseCase(ur repository.UserRepository, pr repository.PostRepository, lr repository.LikeRepository, bmr repository.BookmarkRepository, fr repository.FollowRepository) UserUseCase {
	return &userUseCase{userRepo: ur, postRepo: pr, likeRepo: lr, bmRepo: bmr, followRepo: fr}
}

func (u *userUseCase) GetProfile(ctx context.Context, userID, requestingUserID string) (*domain.User, error) {
	usr, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	followers, err := u.followRepo.CountFollowers(ctx, userID)
	if err == nil {
		usr.FollowersCount = followers
	}
	following, err := u.followRepo.CountFollowings(ctx, userID)
	if err == nil {
		usr.FollowingCount = following
	}
	if requestingUserID != "" {
		followed, err := u.followRepo.IsFollowing(ctx, requestingUserID, userID)
		if err == nil {
			usr.FollowedByMe = followed
		}
	}
	return usr, nil
}

func (u *userUseCase) GetUserPosts(ctx context.Context, userID, requestingUserID string, cursor string, limit int) ([]*domain.Post, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	posts, err := u.postRepo.ListByUser(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	for _, p := range posts {
		count, _ := u.likeRepo.CountByPostID(ctx, p.ID)
		p.LikeCount = count
		bmCount, _ := u.bmRepo.CountByPostID(ctx, p.ID)
		p.BookmarkCount = bmCount
		if requestingUserID != "" {
			liked, _ := u.likeRepo.IsLikedByUser(ctx, p.ID, requestingUserID)
			p.LikedByMe = liked
			bmed, _ := u.bmRepo.IsBookmarkedByUser(ctx, p.ID, requestingUserID)
			p.BookmarkedByMe = bmed
		}
	}
	return posts, nil
}
