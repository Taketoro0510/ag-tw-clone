package usecase_test

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// --- Mocks ---

type mockFirebaseAuth struct {
	verifyIDTokenFn func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *mockFirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.verifyIDTokenFn(ctx, idToken)
}

type mockJWTService struct {
	generateTokenFn func(userID string) (string, error)
}

func (m *mockJWTService) GenerateToken(userID string) (string, error) {
	return m.generateTokenFn(userID)
}

func (m *mockJWTService) VerifyToken(tokenString string) (string, error) {
	return "", nil
}

type mockUserRepository struct {
	findByFirebaseUIDFn func(ctx context.Context, uid string) (*domain.User, error)
	createFn            func(ctx context.Context, user *domain.User) error
	updateFn            func(ctx context.Context, user *domain.User) error
}

func (m *mockUserRepository) FindByFirebaseUID(ctx context.Context, uid string) (*domain.User, error) {
	return m.findByFirebaseUIDFn(ctx, uid)
}
func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.createFn(ctx, user)
}
func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.updateFn(ctx, user)
}

// --- Tests ---

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("success: existing user", func(t *testing.T) {
		mockFA := &mockFirebaseAuth{
			verifyIDTokenFn: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{UID: "firebase-uid-123", Claims: map[string]interface{}{
					"email": "test@example.com",
					"name":  "Test User",
				}}, nil
			},
		}
		mockUserRepo := &mockUserRepository{
			findByFirebaseUIDFn: func(ctx context.Context, uid string) (*domain.User, error) {
				return &domain.User{ID: "user-123", FirebaseUID: "firebase-uid-123"}, nil
			},
			updateFn: func(ctx context.Context, user *domain.User) error {
				return nil
			},
		}
		mockJWT := &mockJWTService{
			generateTokenFn: func(userID string) (string, error) {
				return "jwt-token-123", nil
			},
		}

		u := usecase.NewAuthUseCase(mockFA, mockJWT, mockUserRepo, nil)
		token, err := u.Login(ctx, "valid-id-token")

		assert.NoError(t, err)
		assert.Equal(t, "jwt-token-123", token)
	})

	t.Run("success: new user", func(t *testing.T) {
		mockFA := &mockFirebaseAuth{
			verifyIDTokenFn: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{UID: "firebase-uid-new", Claims: map[string]interface{}{
					"email": "new@example.com",
				}}, nil
			},
		}
		mockUserRepo := &mockUserRepository{
			findByFirebaseUIDFn: func(ctx context.Context, uid string) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			createFn: func(ctx context.Context, user *domain.User) error {
				return nil
			},
		}
		mockJWT := &mockJWTService{
			generateTokenFn: func(userID string) (string, error) {
				return "jwt-token-new", nil
			},
		}

		u := usecase.NewAuthUseCase(mockFA, mockJWT, mockUserRepo, nil)
		token, err := u.Login(ctx, "valid-id-token-new")

		assert.NoError(t, err)
		assert.Equal(t, "jwt-token-new", token)
	})

	t.Run("fail: firebase verify error", func(t *testing.T) {
		mockFA := &mockFirebaseAuth{
			verifyIDTokenFn: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return nil, errors.New("invalid token")
			},
		}

		u := usecase.NewAuthUseCase(mockFA, nil, nil, nil)
		token, err := u.Login(ctx, "invalid-id-token")

		assert.ErrorIs(t, err, domain.ErrUnauthorized)
		assert.Empty(t, token)
	})
}
