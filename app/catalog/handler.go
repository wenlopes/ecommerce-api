package catalog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/product"
	wire_out "github.com/mytheresa/go-hiring-challenge/app/wire/out"
)

type AllProductsResponse struct {
	Products   []wire_out.Product `json:"products"`
	Pagination api.Pagination     `json:"pagination"`
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

// HandleGet handles the HTTP GET request for retrieving products with optional filters and pagination.
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := api.ExtractPagination(r, api.MaxLimit)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Failed to parse pagination parameters")
		return
	}

	filters, err := h.extractProductFilters(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	productsDB, total, err := h.repo.GetAllProducts(offset, limit, filters)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	products := make([]wire_out.Product, len(productsDB))
	for i, p := range productsDB {
		products[i] = wire_out.Product{
			Code:       p.Code,
			Price:      p.Price.InexactFloat64(),
			Categories: make([]wire_out.Category, len(p.Categories)),
		}

		for j, c := range p.Categories {
			products[i].Categories[j] = wire_out.Category{
				Code: c.Code,
				Name: c.Name,
			}
		}
	}

	pagination := api.NewPagination(total, offset, limit)

	response := AllProductsResponse{
		Products:   products,
		Pagination: pagination,
	}

	api.OKResponse(w, http.StatusOK, response)
}

// HandleGetByCode handles the HTTP GET request for retrieving a product by its code.
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Missing product code")
		return
	}

	productDB, err := h.repo.GetProductByCode(code)
	if err != nil {
		if err == product.ErrProductNotFound {
			api.ErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve product data")
		return
	}

	product := wire_out.Product{
		Code:       productDB.Code,
		Price:      productDB.Price.InexactFloat64(),
		Categories: make([]wire_out.Category, len(productDB.Categories)),
		Variants:   make([]wire_out.Variant, len(productDB.Variants)),
	}

	for j, c := range productDB.Categories {
		product.Categories[j] = wire_out.Category{
			Code: c.Code,
			Name: c.Name,
		}
	}

	for k, v := range productDB.Variants {
		product.Variants[k] = wire_out.Variant{
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
