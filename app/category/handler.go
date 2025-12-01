package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Response struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryHandler struct {
	repo Repository
}

// NewCategoryHandler creates a new instance of CategoryHandler.
func NewCategoryHandler(r Repository) *CategoryHandler {
	return &CategoryHandler{
		repo: r,
	}
}

// HandleGet handles the HTTP GET request for retrieving all categories.
func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categoriesDB, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	categories := make([]Category, len(categoriesDB))
	for i, catDB := range categoriesDB {
		categories[i] = Category{
			Code: catDB.Code,
			Name: catDB.Name,
		}
	}

	resp := Response{
		Categories: categories,
	}

	api.OKResponse(w, http.StatusOK, resp)
}

// HandlePost handles the HTTP POST request for creating a new category.
func (h *CategoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var cat Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cat.Code = strings.TrimSpace(cat.Code)
	cat.Name = strings.TrimSpace(cat.Name)

	if cat.Code == "" || cat.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Both code and name are required")
		return
	}

	if err := h.repo.CreateCategory(cat.Code, cat.Name); err != nil {
		if errors.Is(err, ErrCategoryAlreadyExists) {
			api.ErrorResponse(w, http.StatusConflict, "Category already exists")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	api.OKResponse(w, http.StatusCreated, cat)
}
