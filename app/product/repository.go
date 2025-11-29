package product

import "github.com/mytheresa/go-hiring-challenge/models"

type Reader interface {
	GetAllProducts(offset, limit int) ([]models.Product, int64, error)
}

type Repository interface {
	Reader
}
