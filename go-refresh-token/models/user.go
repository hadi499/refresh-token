package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100)" json:"name"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex" json:"email" validate:"required,email"`
	Password     string    `gorm:"type:varchar(255)" json:"password" validate:"required,min=6"`
	RefreshToken string    `gorm:"type:varchar(255);default:''" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate akan dijalankan sebelum data disimpan ke database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}
