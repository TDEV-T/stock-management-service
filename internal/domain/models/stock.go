package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	ImageURL    string    `gorm:"imageurl" json:"imageURL"`
	Description string    `gorm:"column:description" json:"description"`
	CategoryID  uint      `gorm:"column:category_id;not null" json:"categoryId"`
	Category    *Category `json:"category"`
	SKU         string    `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

type ProductDTO struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	ImageURL    string    `gorm:"imageurl" json:"imageURL"`
	Description string    `gorm:"column:description" json:"description"`
	CategoryID  uint      `gorm:"column:category_id;not null" json:"categoryId"`
	Category    *Category `json:"category"`
	SKU         string    `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Quantity    int       `json:"quantity"`
}

type MovementDTO struct {
	Type     string    `json:"type"`
	Quantity int       `json:"quantity"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
	Product  struct {
		Name     string `json:"name"`
		ImageURL string `json:"imageURL"`
	} `json:"product"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}
