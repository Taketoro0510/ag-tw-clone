package domain

import "time"

type Comment struct {
	ID             string
	PostID         string
	AuthorID       string
	Author         *User
	Body           string
	CreatedAt      time.Time
	DeletedAt      *time.Time
	LikeCount      int
	LikedByMe      bool
	BookmarkCount  int
	BookmarkedByMe bool
}
