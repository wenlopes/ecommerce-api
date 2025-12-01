package gorm

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mytheresa/go-hiring-challenge/app/category"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository.
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// GetAllCategories retrieves all categories from the database.
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

// CreateCategory creates a new category in the database.
func (r *CategoryRepository) CreateCategory(code, name string) (models.Category, error) {
	c := models.Category{
		Code: code,
		Name: name,
	}

	err := r.db.Create(&c).Error
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Category{}, category.ErrCategoryAlreadyExists
		}

		return models.Category{}, err
	}

	return c, nil
}
