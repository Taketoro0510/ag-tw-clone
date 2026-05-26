package domain

import "time"

type Bookmark struct {
	ID        string
	UserID    string
	PostID    string
	CreatedAt time.Time
}
