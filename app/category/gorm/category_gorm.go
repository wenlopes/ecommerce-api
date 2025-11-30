package gorm

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.
		Model(&models.Category{}).
		Order("id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
