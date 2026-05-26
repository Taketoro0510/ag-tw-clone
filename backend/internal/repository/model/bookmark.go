package model

import "time"

type Bookmark struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	PostID    string `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}
