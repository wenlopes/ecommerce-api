package category

import "github.com/mytheresa/go-hiring-challenge/models"

type Reader interface {
	GetAllCategories() ([]models.Category, error)
}

type Writer interface {
	CreateCategory(code, name string) (models.Category, error)
}

type Repository interface {
	Reader
	Writer
}
