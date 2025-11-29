package product

import "github.com/mytheresa/go-hiring-challenge/models"

type Filters struct {
	CategoryCode  string
	PriceLessThan *float64
}

type Reader interface {
	GetAllProducts(offset, limit int, filters Filters) ([]models.Product, int64, error)
	GetProductByCode(code string) (models.Product, error)
}

type Repository interface {
	Reader
}
