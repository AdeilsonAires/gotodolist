package repository

import (
	"context"
	"errors"

	"github.com/AdeilsonAires/gotodolist/backend/internal/domain"
	"gorm.io/gorm"
)

type ListRepository interface {
	GetAll(ctx context.Context, userID uint) ([]domain.List, error)
	GetAllActive(ctx context.Context, userID uint) ([]domain.List, error)
	GetByID(ctx context.Context, userID, id uint) (*domain.List, error)
	Create(ctx context.Context, userID uint, list *domain.List) error
	Update(ctx context.Context, userID, id uint, changes domain.ListChanges) error
	Delete(ctx context.Context, userID, id uint) error
}

type listRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) GetAll(ctx context.Context, userID uint) ([]domain.List, error) {
	var lists []domain.List
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&lists).Error
	return lists, err
}

func (r *listRepository) GetAllActive(ctx context.Context, userID uint) ([]domain.List, error) {
	var lists []domain.List
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND completed = ?", userID, false).
		Order("name ASC").
		Find(&lists).Error
	return lists, err
}

func (r *listRepository) GetByID(ctx context.Context, userID, id uint) (*domain.List, error) {
	var list domain.List
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Where("user_id = ? AND id = ?", userID, id).
		First(&list).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &list, nil
}

func (r *listRepository) Create(ctx context.Context, userID uint, list *domain.List) error {
	list.UserID = userID // o dono vem do chamador, nunca do cliente
	list.Completed = false // toda lista nasce ativa, mesmo sem itens
	list.CompletedAt = nil

	err := r.db.WithContext(ctx).Create(list).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrListNameExists
	}
	return err
}

func (r *listRepository) Update(ctx context.Context, userID, id uint, changes domain.ListChanges) error {
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

	// Nenhum campo enviado: não há nada a gravar.
	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&domain.List{}).
		Where("user_id = ? AND id = ?", userID, id).
		Updates(updates)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return domain.ErrListNameExists
	}
	return checkAffected(result)
}

func (r *listRepository) Delete(ctx context.Context, userID, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Apaga a lista primeiro: o filtro por user_id garante o ownership.
		result := tx.
			Where("user_id = ? AND id = ?", userID, id).
			Delete(&domain.List{})
		if err := checkAffected(result); err != nil {
			return err
		}

		// Como é soft delete, o CASCADE do banco não dispara.
		// Apagamos os itens manualmente para não deixar registros órfãos.
		return tx.
			Where("list_id = ?", id).
			Delete(&domain.Item{}).Error
	})
}