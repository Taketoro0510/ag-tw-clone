package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koitake1/cloudcode-sns/backend/internal/auth"
	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(ctx context.Context, idToken string) (string, error)
}

type authUseCase struct {
	firebaseAuth auth.FirebaseAuth
	jwtService   auth.JWTService
	userRepo     repository.UserRepository
	db           *gorm.DB
}

func NewAuthUseCase(fa auth.FirebaseAuth, jwt auth.JWTService, userRepo repository.UserRepository, db *gorm.DB) AuthUseCase {
	return &authUseCase{firebaseAuth: fa, jwtService: jwt, userRepo: userRepo, db: db}
}

func (u *authUseCase) Login(ctx context.Context, idToken string) (string, error) {
	token, err := u.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	user, err := u.userRepo.FindByFirebaseUID(ctx, token.UID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Create user
			id, _ := uuid.NewV7()
			email, _ := token.Claims["email"].(string)
			name, _ := token.Claims["name"].(string)
			if name == "" {
				name = "User"
			}
			var avatarURL *string
			if pic, ok := token.Claims["picture"].(string); ok {
				avatarURL = &pic
			}

			user = &domain.User{
				ID:          id.String(),
				FirebaseUID: token.UID,
				Email:       email,
				DisplayName: name,
				AvatarURL:   avatarURL,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := u.userRepo.Create(ctx, user); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		// Update user info if needed
		email, _ := token.Claims["email"].(string)
		name, _ := token.Claims["name"].(string)
		var avatarURL *string
		if pic, ok := token.Claims["picture"].(string); ok {
			avatarURL = &pic
		}
		user.Email = email
		if name != "" {
			user.DisplayName = name
		}
		if avatarURL != nil {
			user.AvatarURL = avatarURL
		}
		user.UpdatedAt = time.Now()
		u.userRepo.Update(ctx, user)
	}

	jwtStr, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return jwtStr, nil
}
