package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrPostNotFound        = errors.New("post not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrValidation          = errors.New("validation failed")
	ErrConflict            = errors.New("conflict")
	ErrInternalServerError = errors.New("internal server error")
)
