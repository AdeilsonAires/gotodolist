package domain

import (
	"time"

	"gorm.io/gorm"
)

type List struct {
	ID          uint           `json:"id" gorm:"primaryKey;not null"`
	Name        string         `json:"name" validate:"required,min=3,max=20" gorm:"not null"`
	Items       Item           `json:"items" gorm:"foreignKey:ListID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
