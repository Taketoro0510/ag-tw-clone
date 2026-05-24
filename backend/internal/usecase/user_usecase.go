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
	userRepo repository.UserRepository
	postRepo repository.PostRepository
	likeRepo repository.LikeRepository
}

func NewUserUseCase(ur repository.UserRepository, pr repository.PostRepository, lr repository.LikeRepository) UserUseCase {
	return &userUseCase{userRepo: ur, postRepo: pr, likeRepo: lr}
}

func (u *userUseCase) GetProfile(ctx context.Context, userID, requestingUserID string) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
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
		if requestingUserID != "" {
			liked, _ := u.likeRepo.IsLikedByUser(ctx, p.ID, requestingUserID)
			p.LikedByMe = liked
		}
	}
	return posts, nil
}
