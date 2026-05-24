package domain

import "time"

type Post struct {
	ID         string
	AuthorID   string
	Body       string
	MediaType  *string
	MediaPath  *string
	MediaURL   *string
	CreatedAt  time.Time
	DeletedAt  *time.Time
	LikeCount  int
	LikedByMe  bool
}
