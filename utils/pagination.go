package utils

import "gorm.io/gorm"

type PaginationRequest struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func (r PaginationRequest) ComputeMeta(total int64) PaginationMeta {
	totalPage := int(total) / r.Limit
	if int(total)%r.Limit != 0 {
		totalPage++
	}
	if totalPage == 0 {
		totalPage = 1
	}
	return PaginationMeta{
		Page:      r.Page,
		Limit:     r.Limit,
		Total:     total,
		TotalPage: totalPage,
	}
}

func ApplyPagination(query *gorm.DB, pagination PaginationRequest, total *int64) *gorm.DB {
	if !pagination.IsPaginated() {
		return query
	}

	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(total)

	offset := (pagination.Page - 1) * pagination.Limit
	return query.Limit(pagination.Limit).Offset(offset)
}

// func (r *PaginationRequest) SetDefault() {
// 	if r.Limit > 100 {
// 		r.Limit = 100
// 	}
// }

func (r PaginationRequest) IsPaginated() bool {
	return r.Page > 0 && r.Limit > 0
}
