package dto

import (
	"github.com/pedromspeixoto/users-api/internal/data"
)

type PaginationRequest struct {
	Limit  int               `json:"limit,omitempty"`
	Page   int               `json:"page,omitempty"`
	Sort   string            `json:"sort,omitempty"`
	Filter map[string]string `json:"filter,omitempty"`
	Search map[string]string `json:"search,omitempty"`
}

func NewPaginationRequest(limit, page int, sort string, filter, search map[string]string) (*PaginationRequest, error) {
	res := &PaginationRequest{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		Filter: filter,
		Search: search,
	}
	return res, nil
}

func ModelFromPaginationRequest(p *PaginationRequest) *data.Pagination {
	model := &data.Pagination{
		Limit:  p.Limit,
		Page:   p.Page,
		Sort:   p.Sort,
		Filter: p.Filter,
	}
	return model
}

type PaginationResponse struct {
	CurrentPage int         `json:"current_page,omitempty"`
	TotalRows   int64       `json:"total_rows"`
	TotalPages  int         `json:"total_pages"`
	Data        interface{} `json:"data"`
}

func NewPaginationResponse(p *data.Pagination) *PaginationResponse {
	res := &PaginationResponse{
		CurrentPage: p.GetPage(),
		TotalRows:   p.TotalRows,
		TotalPages:  p.TotalPages,
		Data:        p.Data,
	}
	return res
}
