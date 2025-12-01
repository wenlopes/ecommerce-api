package gorm

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/app/product"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

// NewProductsRepository creates a new instance of ProductsRepository.
func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetAllProducts retrieves products from the database with optional filters, pagination, and returns the total count.
func (r *ProductsRepository) GetAllProducts(offset, limit int, filters product.Filters) ([]models.Product, int64, error) {
	var (
		products []models.Product
		total    int64
	)

	baseQuery := r.db.Model(&models.Product{})

	if filters.CategoryCode != "" {
		baseQuery = baseQuery.
			Joins("JOIN product_categories pc ON pc.product_id = products.id").
			Joins("JOIN categories c ON c.id = pc.category_id").
			Where("c.code = ?", filters.CategoryCode)
	}

	if filters.PriceLessThan != nil {
		baseQuery = baseQuery.Where("products.price < ?", *filters.PriceLessThan)
	}

	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.
		Distinct("products.id").
		Count(&total).Error; err != nil {
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

// GetProductByCode retrieves a product by its code, including its variants and categories.
func (r *ProductsRepository) GetProductByCode(code string) (models.Product, error) {
	var p models.Product

	err := r.db.
		Preload("Variants").
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code", "name")
		}).
		Where("code = ?", code).
		First(&p).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p, product.ErrProductNotFound
		}
		return p, err
	}

	return p, nil
}
