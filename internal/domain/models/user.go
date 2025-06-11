package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	LastLoginAt time.Time
}

type UserDTO struct {
	Username    string `gorm:"uniqueIndex;not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	LastLoginAt time.Time
}

func (user *User) ToDTO() *UserDTO {
	return &UserDTO{
		Username:    user.Username,
		Email:       user.Email,
		LastLoginAt: user.LastLoginAt,
	}
}
