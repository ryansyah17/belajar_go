package domain

import "math"

type PaginationParams struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

func (p *PaginationParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func CalculateTotalPage(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}