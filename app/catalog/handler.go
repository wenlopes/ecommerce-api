package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/product"
)

type Response struct {
	Products   []Product      `json:"products"`
	Pagination api.Pagination `json:"pagination"`
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

func (h *CatalogHandler) extractProductFilters(r *http.Request) (product.Filters, error) {
	query := r.URL.Query()
	var filters product.Filters

	if categoryCode := query.Get("category"); categoryCode != "" {
		filters.CategoryCode = categoryCode
	}

	if priceLtStr := query.Get("price_less_than"); priceLtStr != "" {
		priceLt, err := strconv.ParseFloat(priceLtStr, 64)
		if err != nil || priceLt < 0 {
			return filters, fmt.Errorf("invalid price_less_than parameter")
		}
		filters.PriceLessThan = &priceLt
	}

	return filters, nil
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := api.ExtractPagination(r, api.MaxLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filters, err := h.extractProductFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productsDB, total, err := h.repo.GetAllProducts(offset, limit, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(productsDB))
	for i, p := range productsDB {
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

	pagination := api.NewPagination(total, offset, limit)

	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products:   products,
		Pagination: pagination,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
