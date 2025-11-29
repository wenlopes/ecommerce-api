package gorm

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int) ([]models.Product, int64, error) {
	var (
		products []models.Product
		total    int64
	)

	baseQuery := r.db.Model(&models.Product{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.
		Preload("Variants").
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name")
		}).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
