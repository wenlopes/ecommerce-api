package product

import "github.com/mytheresa/go-hiring-challenge/models"

type Reader interface {
	GetAllProducts() ([]models.Product, error)
}

type Repository interface {
	Reader
}
