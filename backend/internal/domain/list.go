package domain

import (
	"time"

	"gorm.io/gorm"
)

type List struct {
	ID          uint           `json:"id" gorm:"primaryKey;not null"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_list_name,where:deleted_at IS NULL"`
	Name        string         `json:"name" validate:"required,min=3,max=20" gorm:"not null;uniqueIndex:idx_user_list_name,where:deleted_at IS NULL"`
	Icon        string         `json:"icon" gorm:"size:50;not null;default:list"`
	Color       string         `json:"color" gorm:"size:7;not null;default:#009966"`
	Items       []Item         `json:"items,omitempty" gorm:"foreignKey:ListID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// ListChanges: nil = não alterar o campo.
type ListChanges struct {
	Name  *string
	Icon  *string
	Color *string
}

func (c ListChanges) IsEmpty() bool {
	return c.Name == nil && c.Icon == nil && c.Color == nil
}