package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;not null"`
	Name      string         `json:"name" validate:"required,min=3,max=20" gorm:"not null"`
	Email     string         `json:"email" validate:"required,min=3,max=100,email" gorm:"not null;unique"`
	Icon      string         `json:"icon" gorm:"size:50;not null;default:user"`
	Color     string         `json:"color" gorm:"size:7;not null;default:#009966"`
	Password  string         `json:"-" validate:"required,min=6,max=20" gorm:"not null"`
	Lists     []List         `json:"lists,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserChanges struct {
	Name     *string
	Email    *string
	Password *string
	Icon     *string
	Color    *string
}

func (c UserChanges) IsEmpty() bool {
	return c.Name == nil && c.Icon == nil && c.Color == nil && c.Email == nil && c.Password == nil
}
