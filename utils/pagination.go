package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"-"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationFromQuery extracts pagination parameters from query string
func GetPaginationFromQuery(c *gin.Context) PaginationRequest {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	offset := (page - 1) * pageSize

	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// BuildPaginationResponse creates a paginated response
func BuildPaginationResponse(data interface{}, total int64, pagination PaginationRequest) PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return PaginationResponse{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}

// GetPaginationMeta creates pagination metadata
func GetPaginationMeta(total int64, pagination PaginationRequest) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return PaginationMeta{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
