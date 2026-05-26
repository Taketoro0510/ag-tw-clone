package model

import "time"

type Post struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	AuthorID  string `gorm:"type:uuid;not null"`
	Body      string `gorm:"not null"`
	MediaType *string
	MediaPath *string
	MediaURL  *string
	CreatedAt time.Time
	DeletedAt *time.Time
	Author    User   `gorm:"foreignKey:AuthorID"`
}
