package domain

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        uint           `json:"id" gorm:"primaryKey;not null"`
	ListID    uint           `json:"list_id" gorm:"not null;uniqueIndex:idx_list_item_name,where:deleted_at IS NULL"`
	Name      string         `json:"name" validate:"required,min=3,max=25" gorm:"not null;uniqueIndex:idx_list_item_name,where:deleted_at IS NULL"`
	Completed bool           `json:"completed" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
