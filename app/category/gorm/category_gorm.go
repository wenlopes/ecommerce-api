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

func (r *CategoryRepository) CreateCategory(code, name string) error {
	cat := models.Category{
		Code: code,
		Name: name,
	}

	if err := r.db.Create(&cat).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return category.ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}
