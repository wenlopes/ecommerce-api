package catalog

import (
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
	Variants   []Variant  `json:"variants,omitempty"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
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
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	filters, err := h.extractProductFilters(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	productsDB, total, err := h.repo.GetAllProducts(offset, limit, filters)
	if err != nil {
		// TODO: Log internal error
		// TODO: Create a generic error message
		api.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
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

	response := Response{
		Products:   products,
		Pagination: pagination,
	}

	api.OKResponse(w, http.StatusOK, response)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "missing product code")
		return
	}

	productDB, err := h.repo.GetProductByCode(code)
	if err != nil {
		// TODO: distinguish not found error
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	product := Product{
		Code:       productDB.Code,
		Price:      productDB.Price.InexactFloat64(),
		Categories: make([]Category, len(productDB.Categories)),
		Variants:   make([]Variant, len(productDB.Variants)),
	}

	for j, c := range productDB.Categories {
		product.Categories[j] = Category{
			Code: c.Code,
			Name: c.Name,
		}
	}

	for k, v := range productDB.Variants {
		product.Variants[k] = Variant{
			Name: v.Name,
			SKU:  v.SKU,
		}
		product.Variants[k].Price = v.Price.InexactFloat64()
		if v.Price.IsZero() {
			product.Variants[k].Price = product.Price
		}
	}

	api.OKResponse(w, http.StatusOK, product)
}
