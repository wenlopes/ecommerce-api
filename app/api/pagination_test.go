package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractPagination(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		maxLimit   int
		wantOffset int
		wantLimit  int
		wantErr    string
	}{
		{
			name:       "defaults applied",
			query:      "",
			maxLimit:   MaxLimit,
			wantOffset: 0,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "valid custom values",
			query:      "offset=5&limit=20",
			maxLimit:   MaxLimit,
			wantOffset: 5,
			wantLimit:  20,
		},
		{
			name:    "invalid offset format",
			query:   "offset=foo",
			wantErr: "invalid offset parameter",
		},
		{
			name:    "negative offset",
			query:   "offset=-1",
			wantErr: "invalid offset parameter",
		},
		{
			name:    "invalid limit format",
			query:   "limit=bar",
			wantErr: "invalid limit parameter",
		},
		{
			name:       "limit lower bound enforced",
			query:      "limit=0",
			maxLimit:   MaxLimit,
			wantOffset: 0,
			wantLimit:  1,
		},
		{
			name:       "limit upper bound enforced",
			query:      "limit=500",
			maxLimit:   50,
			wantOffset: 0,
			wantLimit:  50,
		},
		{
			name:       "max limit corrected when invalid",
			query:      "limit=5",
			maxLimit:   0,
			wantOffset: 0,
			wantLimit:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/catalog"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)

			offset, limit, err := ExtractPagination(req, tt.maxLimit)

			if tt.wantErr != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantErr, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantOffset, offset)
			assert.Equal(t, tt.wantLimit, limit)
		})
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		offset   int
		limit    int
		expected Pagination
	}{
		{
			name:   "limit defaults when non-positive",
			total:  0,
			offset: 0,
			limit:  0,
			expected: Pagination{
				Total:      0,
				Limit:      DefaultLimit,
				Offset:     0,
				Page:       1,
				TotalPages: 0,
				HasNext:    false,
				HasPrev:    false,
			},
		},
		{
			name:   "multiple pages with next/prev",
			total:  25,
			offset: 10,
			limit:  5,
			expected: Pagination{
				Total:      25,
				Limit:      5,
				Offset:     10,
				Page:       3,
				TotalPages: 5,
				HasNext:    true,
				HasPrev:    true,
			},
		},
		{
			name:   "no next page when offset+limit >= total",
			total:  9,
			offset: 6,
			limit:  3,
			expected: Pagination{
				Total:      9,
				Limit:      3,
				Offset:     6,
				Page:       3,
				TotalPages: 3,
				HasNext:    false,
				HasPrev:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPagination(tt.total, tt.offset, tt.limit)
			assert.Equal(t, tt.expected, got)
		})
	}
}
