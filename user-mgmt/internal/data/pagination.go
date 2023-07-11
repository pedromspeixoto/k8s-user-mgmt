package data

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit       int
	Page        int
	Sort        string
	Filter      map[string]string
	Search      map[string]string
	SearchQuery string
	TotalRows   int64
	TotalPages  int
	Data        interface{}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at asc"
	}
	return p.Sort
}

func (p *Pagination) GetFilter() map[string]string {
	if len(p.Filter) == 0 {
		p.Filter = map[string]string{}
	}
	return p.Filter
}

func (p *Pagination) GetSearch() map[string]string {
	if len(p.Search) == 0 {
		p.Search = map[string]string{}
	}
	return p.Search
}

func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	// build search query
	if len(p.GetSearch()) > 0 {
		p.GetSearch()
		first := true
		for field, filter := range p.GetSearch() {
			if first {
				p.SearchQuery = fmt.Sprintf("%s LIKE '%%%s%%'", field, filter)
				first = false
			} else {
				p.SearchQuery = fmt.Sprintf("%s OR %s LIKE '%%%s%%'", p.SearchQuery, field, filter)
			}
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort()).Where(p.GetFilter()).Where(p.SearchQuery)
	}
}

func GetTotalPages(rows int64, limit int) int {
	return int(math.Ceil(float64(rows) / float64(limit)))
}
