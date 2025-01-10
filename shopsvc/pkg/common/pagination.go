package common

import (
	"math"

	"strconv"

	"gorm.io/gorm"

	"shopsvc/internal/core/domain"
)

type Pagination struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

func SetMetaPagination(page int, perPage int, totalData int) *Pagination {
	totalPage := totalData / perPage
	// check if there is a remaining data
	if totalData%perPage > 0 {
		totalPage++
	}

	return &Pagination{
		Page:      page,
		PerPage:   perPage,
		TotalPage: totalPage,
		TotalData: totalData,
	}
}

func AssertPaginationPayload(p domain.QueryParam) *Pagination {
	if p.PerPage == "-" {
		p.PerPage = strconv.Itoa(math.MaxInt64)
	}

	perPage, errPerPage := strconv.Atoi(p.PerPage)
	page, errPage := strconv.Atoi(p.Page)

	if errPerPage != nil {
		perPage = 25
	}

	if errPage != nil {
		page = 1
	}

	return &Pagination{
		PerPage: int(perPage),
		Page:    int(page),
	}
}

func Paginate(ctx ServiceContextManager, value interface{}, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var (
		page      int    = ctx.GetPagination().Page
		perPage   int    = ctx.GetPagination().PerPage
		sort      string = ctx.GetSortParam().Sort
		totalRows int64
	)

	// Count total data
	db.Model(value).Count(&totalRows)

	// Set meta pagination to context
	ctx.SetPagination(SetMetaPagination(page, perPage, int(totalRows)))

	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case perPage > 100:
			perPage = 100
		case perPage <= 0:
			perPage = 10
		}

		offset := (page - 1) * perPage

		return db.Offset(offset).Limit(perPage).Order(sort)
	}
}
