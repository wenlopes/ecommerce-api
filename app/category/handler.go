package category

import (
	"encoding/json"
	"net/http"
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
