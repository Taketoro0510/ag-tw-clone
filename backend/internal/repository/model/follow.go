package model

import "time"

type Follow struct {
	FollowerID string    `gorm:"primaryKey;column:follower_id"`
	FolloweeID string    `gorm:"primaryKey;column:followee_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}
