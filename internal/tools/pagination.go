package tools

import (
	"github.com/labstack/echo/v4"
	"math"
	"strconv"
)

//Pagination struct
type Pagination struct {
	TotalItem int64
	PageSize  int32
	Page      int32
	Data      interface{}
}

//PaginationResponse struct
type PaginationResponse struct {
	Message         string      `json:"message"`
	Data            interface{} `json:"data"`
	CurrentPage     int32       `json:"current_page"`
	HasPreviousPage bool        `json:"has_previous_page"`
	HasNextPage     bool        `json:"has_next_page"`
	PreviousPage    int32       `json:"previous_page"`
	NextPage        int32       `json:"next_page"`
	LastPage        int32       `json:"last_page"`
	TotalItem       int64       `json:"total_item"`
}

//PaginationPageAndPageSize represent value of page and pageSize based on query params headers
func PaginationPageAndPageSize(c echo.Context) (page int32, pageSize int32) {
	pageParam := c.QueryParam("page")
	pageSizeParam := c.QueryParam("page_size")

	pageVal := 0
	pageSizeVal := 0

	pageVal, err := strconv.Atoi(pageParam)
	if err != nil {
		pageVal = 1
	}
	pageSizeVal, err = strconv.Atoi(pageSizeParam)
	if err != nil {
		pageSizeVal = 10
	}
	return int32(pageVal), int32(pageSizeVal)
}

//PaginationGetResponse represent return data json response for pagination
func PaginationGetResponse(message string, pagination Pagination) PaginationResponse {
	lastPageFloat := float64(pagination.TotalItem) / float64(pagination.PageSize)
	lastPage := int32(math.Ceil(lastPageFloat))
	return PaginationResponse{
		Message:         message,
		Data:            pagination.Data,
		CurrentPage:     pagination.Page,
		HasPreviousPage: pagination.Page > 1,
		HasNextPage:     pagination.Page*pagination.PageSize < int32(pagination.TotalItem),
		PreviousPage:    pagination.Page - 1,
		NextPage:        pagination.Page + 1,
		LastPage:        lastPage,
		TotalItem:       pagination.TotalItem,
	}
}

//PaginationSkip return offset value based on page and pageSize value
func PaginationSkip(page int32, pageSize int32) int32 {
	skip := (page - 1) * pageSize
	return skip
}

//PaginationLimit represent return data for page size on infinite scroll
func PaginationLimit(c echo.Context) int32 {
	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}
	return int32(limit)
}
