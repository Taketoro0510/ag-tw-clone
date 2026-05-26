package usecase

import (
	"context"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
)

type BookmarkUseCase interface {
	BookmarkPost(ctx context.Context, userID, postID string) error
	UnbookmarkPost(ctx context.Context, userID, postID string) error
	ListBookmarks(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error)
}

type bookmarkUseCase struct {
	bmRepo   repository.BookmarkRepository
	likeRepo repository.LikeRepository
}

func NewBookmarkUseCase(bmr repository.BookmarkRepository, lr repository.LikeRepository) BookmarkUseCase {
	return &bookmarkUseCase{bmRepo: bmr, likeRepo: lr}
}

func (u *bookmarkUseCase) BookmarkPost(ctx context.Context, userID, postID string) error {
	return u.bmRepo.Bookmark(ctx, userID, postID)
}

func (u *bookmarkUseCase) UnbookmarkPost(ctx context.Context, userID, postID string) error {
	return u.bmRepo.Unbookmark(ctx, userID, postID)
}

func (u *bookmarkUseCase) ListBookmarks(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error) {
	posts, err := u.bmRepo.ListBookmarkedPosts(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	for _, p := range posts {
		count, _ := u.likeRepo.CountByPostID(ctx, p.ID)
		p.LikeCount = count
		if userID != "" {
			liked, _ := u.likeRepo.IsLikedByUser(ctx, p.ID, userID)
			p.LikedByMe = liked
			bmed, _ := u.bmRepo.IsBookmarkedByUser(ctx, p.ID, userID)
			p.BookmarkedByMe = bmed
		}
	}
	return posts, nil
}
