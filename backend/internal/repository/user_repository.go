package repository

import (
	"context"
	"errors"

	"github.com/AdeilsonAires/gotodolist/backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, id uint, changes domain.UserChanges) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrEmailAlreadyExists
	}
	return err
}

func (r *userRepository) Update(ctx context.Context, id uint, changes domain.UserChanges) error {
	updates := map[string]any{}
	if changes.Name != nil {
		updates["name"] = *changes.Name
	}
	if changes.Icon != nil {
		updates["icon"] = *changes.Icon
	}
	if changes.Color != nil {
		updates["color"] = *changes.Color
	}
	if changes.Email != nil {
		updates["email"] = *changes.Email
	}
	if changes.Password != nil {
		updates["password"] = *changes.Password
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(updates)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return domain.ErrEmailAlreadyExists
	}
	return checkAffected(result)
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&domain.User{})
		if err := checkAffected(result); err != nil {
			return err
		}

		// Soft delete não dispara o CASCADE do banco, então limpamos à mão.
		// Ordem importa: itens ANTES das listas, porque a subquery só enxerga
		// listas que ainda não foram apagadas (filtro de deleted_at).
		ownedLists := tx.Model(&domain.List{}).
			Select("id").
			Where("user_id = ?", id)

		if err := tx.Where("list_id IN (?)", ownedLists).
			Delete(&domain.Item{}).Error; err != nil {
			return err
		}

		return tx.Where("user_id = ?", id).
			Delete(&domain.List{}).Error
	})
}