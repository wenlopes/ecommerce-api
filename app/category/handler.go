package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	wire_in "github.com/mytheresa/go-hiring-challenge/app/wire/in"
	wire_out "github.com/mytheresa/go-hiring-challenge/app/wire/out"
)

type AllCategoriesResponse struct {
	Categories []wire_out.Category `json:"categories"`
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

	categories := make([]wire_out.Category, len(categoriesDB))
	for i, catDB := range categoriesDB {
		categories[i] = wire_out.Category{
			Code: catDB.Code,
			Name: catDB.Name,
		}
	}

	resp := AllCategoriesResponse{
		Categories: categories,
	}

	api.OKResponse(w, http.StatusOK, resp)
}

// HandlePost handles the HTTP POST request for creating a new category.
func (h *CategoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var req wire_in.Category

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Both code and name are required")
		return
	}

	created, err := h.repo.CreateCategory(req.Code, req.Name)
	if err != nil {
		if errors.Is(err, ErrCategoryAlreadyExists) {
			api.ErrorResponse(w, http.StatusConflict, "Category already exists")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	category := wire_out.Category{
		Code: created.Code,
		Name: created.Name,
	}

	api.OKResponse(w, http.StatusCreated, category)
}
