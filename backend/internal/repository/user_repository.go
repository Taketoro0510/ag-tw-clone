package repository

import (
	"context"
	"errors"

	"github.com/koitake1/cloudcode-sns/backend/internal/domain"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByFirebaseUID(ctx context.Context, uid string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) toDomain(m *model.User) *domain.User {
	return &domain.User{
		ID:          m.ID,
		FirebaseUID: m.FirebaseUID,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		AvatarURL:   m.AvatarURL,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *userRepository) toModel(d *domain.User) *model.User {
	return &model.User{
		ID:          d.ID,
		FirebaseUID: d.FirebaseUID,
		Email:       d.Email,
		DisplayName: d.DisplayName,
		AvatarURL:   d.AvatarURL,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func (r *userRepository) FindByFirebaseUID(ctx context.Context, uid string) (*domain.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("firebase_uid = ?", uid).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	m := r.toModel(user)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	m := r.toModel(user)
	return r.db.WithContext(ctx).Save(m).Error
}
