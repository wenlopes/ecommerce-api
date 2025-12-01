package api

import "net/http"

type Handlers struct {
	Catalog  CatalogRoutes
	Category CategoryRoutes
}

// CatalogRoutes exposes the catalog endpoints used by the HTTP router.
type CatalogRoutes interface {
	HandleGet(http.ResponseWriter, *http.Request)
	HandleGetByCode(http.ResponseWriter, *http.Request)
}

// CategoryRoutes exposes the category endpoints used by the HTTP router.
type CategoryRoutes interface {
	HandleGet(http.ResponseWriter, *http.Request)
	HandlePost(http.ResponseWriter, *http.Request)
}

// NewMuxRouter wires the HTTP routes for the API and returns a configured ServeMux.
func NewMuxRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	// Catalog routes
	mux.HandleFunc("GET /catalog", h.Catalog.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", h.Catalog.HandleGetByCode)

	// Category routes
	mux.HandleFunc("GET /category", h.Category.HandleGet)
	mux.HandleFunc("POST /category", h.Category.HandlePost)
	return mux
}
