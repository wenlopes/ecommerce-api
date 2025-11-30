package category

import "github.com/mytheresa/go-hiring-challenge/models"

type Reader interface {
	GetAllCategories() ([]models.Category, error)
}

type Repository interface {
	Reader
}
