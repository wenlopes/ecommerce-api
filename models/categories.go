package models

// Category represents a product category used to group products together.
// Each category has a unique code and human-readable name.
type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Code     string    `gorm:"uniqueIndex;not null"`
	Name     string    `gorm:"not null"`
	Products []Product `gorm:"many2many:product_categories"`
}

func (c *Category) TableName() string {
	return "categories"
}
