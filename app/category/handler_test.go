package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	categorymock "github.com/mytheresa/go-hiring-challenge/app/category/mock"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCategoryHandler_HandleGet(t *testing.T) {
	t.Run("RepositoryError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := categorymock.NewMockRepository(ctrl)
		repoErr := errors.New("db down")
		repo.EXPECT().
			GetAllCategories().
			Return(nil, repoErr)

		handler := NewCategoryHandler(repo)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.JSONEq(t, `{"error":"Failed to retrieve categories"}`, recorder.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		categories := []models.Category{
			{Code: "cat-001", Name: "T-Shirts"},
			{Code: "cat-002", Name: "Shoes"},
		}

		repo := categorymock.NewMockRepository(ctrl)
		repo.EXPECT().
			GetAllCategories().
			Return(categories, nil)

		handler := NewCategoryHandler(repo)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var resp Response
		if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, []Category{
			{Code: "cat-001", Name: "T-Shirts"},
			{Code: "cat-002", Name: "Shoes"},
		}, resp.Categories)
	})
}

func TestCategoryHandler_HandlePost(t *testing.T) {
	t.Run("InvalidPayload", func(t *testing.T) {
		handler := &CategoryHandler{}
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("not-json"))
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"Invalid request payload"}`, recorder.Body.String())
	})

	t.Run("MissingFields", func(t *testing.T) {
		handler := &CategoryHandler{}

		body := bytes.NewBufferString(`{"code":"   ","name":"   "}`)
		req := httptest.NewRequest(http.MethodPost, "/categories", body)
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"Both code and name are required"}`, recorder.Body.String())
	})

	t.Run("CategoryAlreadyExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := categorymock.NewMockRepository(ctrl)
		repo.EXPECT().
			CreateCategory("cat-001", "Tops").
			Return(ErrCategoryAlreadyExists)

		handler := NewCategoryHandler(repo)

		body := bytes.NewBufferString(`{"code":"cat-001","name":"Tops"}`)
		req := httptest.NewRequest(http.MethodPost, "/categories", body)
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusConflict, recorder.Code)
		assert.JSONEq(t, `{"error":"Category already exists"}`, recorder.Body.String())
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := categorymock.NewMockRepository(ctrl)
		repoErr := errors.New("db down")
		repo.EXPECT().
			CreateCategory("cat-009", "Accessories").
			Return(repoErr)

		handler := NewCategoryHandler(repo)

		body := bytes.NewBufferString(`{"code":"cat-009","name":"Accessories"}`)
		req := httptest.NewRequest(http.MethodPost, "/categories", body)
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.JSONEq(t, `{"error":"Failed to create category"}`, recorder.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		repo := categorymock.NewMockRepository(ctrl)
		repo.EXPECT().
			CreateCategory("cat-555", "Dresses").
			Return(nil)

		handler := NewCategoryHandler(repo)

		body := bytes.NewBufferString(`{"code":" cat-555 ","name":" Dresses "}`)
		req := httptest.NewRequest(http.MethodPost, "/categories", body)
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var resp Category
		if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, Category{Code: "cat-555", Name: "Dresses"}, resp)
	})
}
