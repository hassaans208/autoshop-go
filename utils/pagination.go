package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageParams struct {
	Page  int
	Limit int
	Q     string
}

func ParsePage(c *gin.Context) PageParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return PageParams{Page: page, Limit: limit, Q: c.Query("q")}
}

func Offset(p PageParams) int {
	return (p.Page - 1) * p.Limit
}

type PagedResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func Paged(data any, total int64, p PageParams) PagedResponse {
	return PagedResponse{Data: data, Total: total, Page: p.Page, Limit: p.Limit}
}
