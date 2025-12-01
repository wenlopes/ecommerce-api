package models

import "time"

// ProductCategory represents the join table that links products and categories.
type ProductCategory struct {
	ID         uint `gorm:"primaryKey"`
	ProductID  uint `gorm:"not null"`
	CategoryID uint `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (pc *ProductCategory) TableName() string {
	return "product_categories"
}
