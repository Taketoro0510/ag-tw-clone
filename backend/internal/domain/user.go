package domain

import "time"

type User struct {
	ID             string
	FirebaseUID    string
	Email          string
	DisplayName    string
	AvatarURL      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FollowersCount int
	FollowingCount int
	FollowedByMe   bool
}
