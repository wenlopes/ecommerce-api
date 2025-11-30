package category

import (
	"encoding/json"
	"net/http"
	"strings"
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

func NewCategoryHandler(r Repository) *CategoryHandler {
	return &CategoryHandler{
		repo: r,
	}
}

func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categoriesDB, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CategoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var cat Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	cat.Code = strings.TrimSpace(cat.Code)
	cat.Name = strings.TrimSpace(cat.Name)

	if cat.Code == "" || cat.Name == "" {
		http.Error(w, "Both code and name are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateCategory(cat.Code, cat.Name); err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
