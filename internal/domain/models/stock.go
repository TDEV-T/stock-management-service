package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Category    string `gorm:"not null"`
	SKU         string `gorm:"uniqueIndex;not null"`
}

type Stock struct {
	gorm.Model
	ProductID uint
	Product   Product
	Quantity  int `gorm:"not null"`
}

type StockMovement struct {
	gorm.Model
	ProductID uint
	Product   Product
	UserID    uint
	User      User
	Type      string    `gorm:"not null"` // "import" or "export"
	Quantity  int       `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	Notes     string
}
