package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;not null"`
	Name      string         `json:"name" validate:"required,min=3,max=20" gorm:"not null"`
	Email     string         `json:"email" validate:"required,min=3,max=25,email" gorm:"not null;unique"`
	Icon      string         `json:"icon" validate:"omitempty,max=255" gorm:"size:255"`
	Color     string         `json:"color" validate:"omitempty,hexcolor" gorm:"size:7"`
	Password  string         `json:"-" validate:"required,min=6,max=20" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
