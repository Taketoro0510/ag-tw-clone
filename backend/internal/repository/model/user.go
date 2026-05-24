package model

import "time"

type User struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	FirebaseUID string `gorm:"uniqueIndex;not null"`
	Email       string `gorm:"not null"`
	DisplayName string `gorm:"not null"`
	AvatarURL   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
