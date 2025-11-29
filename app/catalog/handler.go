package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/product"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code       string     `json:"code"`
	Price      float64    `json:"price"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CatalogHandler struct {
	repo product.Repository
}

func NewCatalogHandler(r product.Repository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:       p.Code,
			Price:      p.Price.InexactFloat64(),
			Categories: make([]Category, len(p.Categories)),
		}
		for j, c := range p.Categories {
			products[i].Categories[j] = Category{
				Code: c.Code,
				Name: c.Name,
			}
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
