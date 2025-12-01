package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

type Pagination struct {
	Total      int64 `json:"total"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Page       int   `json:"page"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ExtractPagination extracts pagination parameters from the HTTP request.
func ExtractPagination(r *http.Request, maxLimit int) (offset, limit int, err error) {
	query := r.URL.Query()

	offsetParam := query.Get("offset")
	if offsetParam != "" {
		offsetConverted, convErr := strconv.Atoi(offsetParam)
		if convErr != nil || offsetConverted < 0 {
			return 0, 0, fmt.Errorf("invalid offset parameter")
		}
		offset = offsetConverted
	}

	if maxLimit < 1 {
		maxLimit = 1
	}

	limit = DefaultLimit
	limitParam := query.Get("limit")
	if limitParam != "" {
		limitConverted, convErr := strconv.Atoi(limitParam)
		if convErr != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}

		if limitConverted < 1 {
			limitConverted = 1
		}
		if limitConverted > maxLimit {
			limitConverted = maxLimit
		}
		limit = limitConverted
	}

	return offset, limit, nil
}

// NewPagination creates a new Pagination struct based on the total number of items, offset, and limit.
func NewPagination(total int64, offset, limit int) Pagination {
	if limit < 1 {
		limit = DefaultLimit
	}

	page := (offset / limit) + 1
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := offset+limit < int(total)
	hasPrev := offset > 0

	return Pagination{
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}
