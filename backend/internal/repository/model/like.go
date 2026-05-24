package model

import "time"

type Like struct {
	PostID    string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
}
