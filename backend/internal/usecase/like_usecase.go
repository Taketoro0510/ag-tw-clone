package usecase

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
)

type LikeUseCase interface {
	LikePost(ctx context.Context, userID, postID string) error
	UnlikePost(ctx context.Context, userID, postID string) error
}

type likeUseCase struct {
	likeRepo repository.LikeRepository
}

func NewLikeUseCase(lr repository.LikeRepository) LikeUseCase {
	return &likeUseCase{likeRepo: lr}
}

func (u *likeUseCase) LikePost(ctx context.Context, userID, postID string) error {
	return u.likeRepo.Create(ctx, postID, userID)
}

func (u *likeUseCase) UnlikePost(ctx context.Context, userID, postID string) error {
	return u.likeRepo.Delete(ctx, postID, userID)
}
