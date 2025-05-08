package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams adalah parameter standar untuk pagination
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationResult adalah hasil dari pagination
type PaginationResult struct {
	Total       int         `json:"total"`
	Page        int         `json:"page"`
	Limit       int         `json:"limit"`
	TotalPages  int         `json:"total_pages"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	Data        interface{} `json:"data"`
}

// DefaultPage adalah nilai default untuk page
const DefaultPage = 1

// DefaultLimit adalah nilai default untuk limit
const DefaultLimit = 10

// GetPaginationParams mengambil parameter pagination dari request
func GetPaginationParams(c *gin.Context) PaginationParams {
	params := PaginationParams{
		Page:  DefaultPage,
		Limit: DefaultLimit,
	}

	// Get page parameter
	pageStr := c.DefaultQuery("page", "")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			params.Page = page
		}
	}

	// Get limit parameter
	limitStr := c.DefaultQuery("limit", "")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	return params
}

// CreatePaginationResult membuat hasil pagination
func CreatePaginationResult(total int, params PaginationParams, data interface{}) PaginationResult {
	totalPages := (total + params.Limit - 1) / params.Limit

	return PaginationResult{
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrevious: params.Page > 1,
		Data:        data,
	}
}
