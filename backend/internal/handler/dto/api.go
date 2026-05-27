package dto

import "time"

type User struct {
	ID             string    `json:"id"`
	FirebaseUID    string    `json:"firebaseUid"`
	Email          string    `json:"email"`
	DisplayName    string    `json:"displayName"`
	AvatarURL      *string   `json:"avatarUrl"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	FollowersCount int       `json:"followersCount"`
	FollowingCount int       `json:"followingCount"`
	FollowedByMe   bool      `json:"followedByMe"`
}

type Post struct {
	ID             string    `json:"id"`
	AuthorID       string    `json:"authorId"`
	Author         *User     `json:"author,omitempty"`
	Body           string    `json:"body"`
	MediaType      *string   `json:"mediaType" enums:"image,video"`
	MediaPath      *string   `json:"mediaPath"`
	MediaURL       *string   `json:"mediaUrl"`
	LikeCount      int       `json:"likeCount"`
	LikedByMe      bool      `json:"likedByMe"`
	BookmarkCount  int       `json:"bookmarkCount"`
	BookmarkedByMe bool      `json:"bookmarkedByMe"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Comment struct {
	ID             string    `json:"id"`
	PostID         string    `json:"postId"`
	AuthorID       string    `json:"authorId"`
	Author         *User     `json:"author,omitempty"`
	Body           string    `json:"body"`
	LikeCount      int       `json:"likeCount"`
	LikedByMe      bool      `json:"likedByMe"`
	BookmarkCount  int       `json:"bookmarkCount"`
	BookmarkedByMe bool      `json:"bookmarkedByMe"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateSessionRequest struct {
	IdToken string `json:"idToken"`
}

type CreateSessionResponse struct {
	Token string `json:"token"`
}

type CreatePostRequest struct {
	Body      string  `json:"body"`
	MediaType *string `json:"mediaType" enums:"image,video"`
	MediaPath *string `json:"mediaPath"`
}

type PaginatedPostsResponse struct {
	Items      []Post  `json:"items"`
	NextCursor *string `json:"nextCursor"`
}

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type PaginatedCommentsResponse struct {
	Items      []Comment `json:"items"`
	NextCursor *string   `json:"nextCursor"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
