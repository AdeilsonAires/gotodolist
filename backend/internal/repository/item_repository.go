package repository

import (
	"context"
	"errors"

	"github.com/AdeilsonAires/gotodolist/backend/internal/domain"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAll(ctx context.Context, listID, userID uint) ([]domain.Item, error)
	GetByID(ctx context.Context, listID, id, userID uint) (*domain.Item, error)
	Rename(ctx context.Context, listID, id, userID uint, newName string) error
	Create(ctx context.Context, item *domain.Item, userID uint) error
	Delete(ctx context.Context, listID, id, userID uint) error
	SetCompleted(ctx context.Context, listID, id, userID uint, completed bool) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) GetAll(ctx context.Context, listID, userID uint) ([]domain.Item, error) {
	var items []domain.Item
	err := scoped(r.db.WithContext(ctx), listID, userID).
		Order("name ASC").
		Find(&items).Error
	return items, err
}

func (r *itemRepository) GetByID(ctx context.Context, listID, id, userID uint) (*domain.Item, error) {
	var item domain.Item
	err := scoped(r.db.WithContext(ctx), listID, userID).
		Where("id = ?", id).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// Rename não altera o status da lista, então não precisa de transação.
func (r *itemRepository) Rename(ctx context.Context, listID, id, userID uint, newName string) error {
	result := scoped(r.db.WithContext(ctx), listID, userID).
		Where("id = ?", id).
		Update("name", newName)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return domain.ErrItemNameExists
	}
	return checkAffected(result)
}

func (r *itemRepository) Create(ctx context.Context, item *domain.Item, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockList(tx, item.ListID, userID); err != nil {
			return err
		}

		item.Completed = false // item novo sempre nasce incompleto

		err := tx.Create(item).Error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrItemNameExists
		}
		if err != nil {
			return err
		}

		return syncListCompletion(tx, item.ListID)
	})
}

func (r *itemRepository) Delete(ctx context.Context, listID, id, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockList(tx, listID, userID); err != nil {
			return err
		}

		result := scoped(tx, listID, userID).
			Where("id = ?", id).
			Delete(&domain.Item{})
		if err := checkAffected(result); err != nil {
			return err
		}

		return syncListCompletion(tx, listID)
	})
}

func (r *itemRepository) SetCompleted(ctx context.Context, listID, id, userID uint, completed bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockList(tx, listID, userID); err != nil {
			return err
		}

		result := scoped(tx, listID, userID).
			Where("id = ?", id).
			Update("completed", completed)
		if err := checkAffected(result); err != nil {
			return err
		}

		return syncListCompletion(tx, listID)
	})
}
