package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/product"
	productmock "github.com/mytheresa/go-hiring-challenge/app/product/mock"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCatalogHandler_HandleGetByCode(t *testing.T) {
	t.Run("MissingCode", func(t *testing.T) {
		handler := &CatalogHandler{}

		req := httptest.NewRequest(http.MethodGet, "/catalog/product", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"Missing product code"}`, recorder.Body.String())
	})

	t.Run("NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := productmock.NewMockRepository(ctrl)
		repo.EXPECT().
			GetProductByCode("ABC").
			Return(models.Product{}, product.ErrProductNotFound)

		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/catalog/product/ABC", nil)
		req.SetPathValue("code", "ABC")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.JSONEq(t, `{"error":"Product not found"}`, recorder.Body.String())
	})

	t.Run("RepositoryError", func(t *testing.T) {
		repoErr := errors.New("database offline")

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := productmock.NewMockRepository(ctrl)
		repo.EXPECT().
			GetProductByCode("XYZ").
			Return(models.Product{}, repoErr)

		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/catalog/product/XYZ", nil)
		req.SetPathValue("code", "XYZ")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.JSONEq(t, `{"error":"`+repoErr.Error()+`"}`, recorder.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		productPrice := decimal.RequireFromString("10.50")
		variantPrice := decimal.RequireFromString("15.00")

		productModel := models.Product{
			Code:  "SKU-001",
			Price: productPrice,
			Categories: []models.Category{
				{Code: "cat-001", Name: "T-Shirts"},
				{Code: "cat-002", Name: "Summer"},
			},
			Variants: []models.Variant{
				{Name: "Small", SKU: "SKU-001-S", Price: decimal.Zero},
				{Name: "Large", SKU: "SKU-001-L", Price: variantPrice},
			},
		}

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := productmock.NewMockRepository(ctrl)
		repo.EXPECT().
			GetProductByCode(productModel.Code).
			Return(productModel, nil)

		handler := NewCatalogHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/catalog/product/SKU-001", nil)
		req.SetPathValue("code", "SKU-001")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", recorder.Code)
		}
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var got Product
		if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// Product
		assert.Equal(t, productModel.Code, got.Code)
		assert.Equal(t, productModel.Price.InexactFloat64(), got.Price)

		// Categories
		if len(got.Categories) != len(productModel.Categories) {
			t.Fatalf("expected %d categories, got %d", len(productModel.Categories), len(got.Categories))
		}
		assert.Equal(t, productModel.Categories[0].Code, got.Categories[0].Code)
		assert.Equal(t, productModel.Categories[0].Name, got.Categories[0].Name)
		assert.Equal(t, productModel.Categories[1].Code, got.Categories[1].Code)
		assert.Equal(t, productModel.Categories[1].Name, got.Categories[1].Name)

		// Variants
		if len(got.Variants) != len(productModel.Variants) {
			t.Fatalf("expected %d variants, got %d", len(productModel.Variants), len(got.Variants))
		}

		// Variant 0 herda o preço do produto
		assert.Equal(t, productModel.Variants[0].Name, got.Variants[0].Name)
		assert.Equal(t, productModel.Variants[0].SKU, got.Variants[0].SKU)
		assert.Equal(t, productModel.Price.InexactFloat64(), got.Variants[0].Price)

		// Variant 1 usa o próprio preço
		assert.Equal(t, productModel.Variants[1].Name, got.Variants[1].Name)
		assert.Equal(t, productModel.Variants[1].SKU, got.Variants[1].SKU)
		assert.Equal(t, productModel.Variants[1].Price.InexactFloat64(), got.Variants[1].Price)
	})
}
