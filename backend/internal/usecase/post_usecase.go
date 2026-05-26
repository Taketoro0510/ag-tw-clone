package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
	"gorm.io/gorm"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, authorID string, body string, mediaType, mediaPath, mediaURL *string) (*domain.Post, error)
	DeletePost(ctx context.Context, userID, postID string) error
	GetTimeline(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error)
	GetPost(ctx context.Context, userID, postID string) (*domain.Post, error)
}

type postUseCase struct {
	postRepo repository.PostRepository
	likeRepo repository.LikeRepository
	bmRepo   repository.BookmarkRepository
	db       *gorm.DB
}

func NewPostUseCase(postRepo repository.PostRepository, likeRepo repository.LikeRepository, bmRepo repository.BookmarkRepository, db *gorm.DB) PostUseCase {
	return &postUseCase{postRepo: postRepo, likeRepo: likeRepo, bmRepo: bmRepo, db: db}
}

func (u *postUseCase) CreatePost(ctx context.Context, authorID string, body string, mediaType, mediaPath, mediaURL *string) (*domain.Post, error) {
	if len(body) > 140 {
		return nil, domain.ErrValidation
	}
	id, _ := uuid.NewV7()
	post := &domain.Post{
		ID:        id.String(),
		AuthorID:  authorID,
		Body:      body,
		MediaType: mediaType,
		MediaPath: mediaPath,
		MediaURL:  mediaURL,
		CreatedAt: time.Now(),
	}
	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (u *postUseCase) DeletePost(ctx context.Context, userID, postID string) error {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return domain.ErrForbidden
	}
	return u.postRepo.Delete(ctx, postID)
}

func (u *postUseCase) GetTimeline(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Post, error) {
	posts, err := u.postRepo.ListGlobal(ctx, cursor, limit)
	if err != nil {
		return nil, err
	}
	for _, p := range posts {
		count, _ := u.likeRepo.CountByPostID(ctx, p.ID)
		p.LikeCount = count
		bmCount, _ := u.bmRepo.CountByPostID(ctx, p.ID)
		p.BookmarkCount = bmCount
		if userID != "" {
			liked, _ := u.likeRepo.IsLikedByUser(ctx, p.ID, userID)
			p.LikedByMe = liked
			bmed, _ := u.bmRepo.IsBookmarkedByUser(ctx, p.ID, userID)
			p.BookmarkedByMe = bmed
		}
	}
	return posts, nil
}

func (u *postUseCase) GetPost(ctx context.Context, userID, postID string) (*domain.Post, error) {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	count, _ := u.likeRepo.CountByPostID(ctx, post.ID)
	post.LikeCount = count
	bmCount, _ := u.bmRepo.CountByPostID(ctx, post.ID)
	post.BookmarkCount = bmCount
	if userID != "" {
		liked, _ := u.likeRepo.IsLikedByUser(ctx, post.ID, userID)
		post.LikedByMe = liked
		bmed, _ := u.bmRepo.IsBookmarkedByUser(ctx, post.ID, userID)
		post.BookmarkedByMe = bmed
	}
	return post, nil
}
