package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds the parsed query parameters for pagination, sorting, and searching.
type PaginationParams struct {
	Page       int
	Limit      int
	Offset     int
	Search     string
	Sort       string
	Order      string
	IsPaginate bool // True if page or limit were provided
}

// GetPaginationParams parses and validates pagination parameters from the Gin context.
// It uses a whitelist to ensure only allowed columns can be used for sorting to prevent SQL injection.
func GetPaginationParams(c *gin.Context, defaultSort string, allowedSortColumns map[string]string) PaginationParams {
	isPaginate := false

	// Parse page
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err == nil && page > 0 {
		isPaginate = true
	} else {
		page = 1
	}

	// Parse limit
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err == nil && limit > 0 {
		isPaginate = true
		if limit > 100 {
			limit = 100 // max limit
		}
	} else {
		limit = 1000 // default large limit if not paginating
	}

	offset := (page - 1) * limit

	// Parse search
	search := c.Query("search")

	// Parse sort
	sortParam := c.Query("sort")
	sortCol, ok := allowedSortColumns[sortParam]
	if !ok || sortParam == "" {
		sortCol = defaultSort
	}

	// Parse order
	orderParam := strings.ToUpper(c.Query("order"))
	if orderParam != "ASC" && orderParam != "DESC" {
		orderParam = "ASC" // default order
	}

	return PaginationParams{
		Page:       page,
		Limit:      limit,
		Offset:     offset,
		Search:     search,
		Sort:       sortCol,
		Order:      orderParam,
		IsPaginate: isPaginate,
	}
}
