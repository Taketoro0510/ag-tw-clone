package model

import "time"

type Comment struct {
	ID        string     `gorm:"primaryKey;column:id"`
	PostID    string     `gorm:"column:post_id"`
	AuthorID  string     `gorm:"column:author_id"`
	Author    User       `gorm:"foreignKey:AuthorID"`
	Body      string     `gorm:"column:body"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

type CommentLike struct {
	CommentID string    `gorm:"primaryKey;column:comment_id"`
	UserID    string    `gorm:"primaryKey;column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type CommentBookmark struct {
	CommentID string    `gorm:"primaryKey;column:comment_id"`
	UserID    string    `gorm:"primaryKey;column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
