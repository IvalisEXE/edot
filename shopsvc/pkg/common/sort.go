package common

import (
	"shopsvc/internal/core/domain"
	"strings"

	"gorm.io/gorm"
)

type SortParam struct {
	Sort string
}

func AssertSortPayload(p domain.QueryParam) *SortParam {
	var param []string

	if p.Sort == "" {
		return &SortParam{}
	}

	sorts := strings.Split(p.Sort, ",")

	for _, sort := range sorts {
		sortType := "ASC" // Asc

		if sort[0] == '-' {
			sortType = "DESC" // Desc
			sort = sort[1:]
		}

		param = append(param, sort+" "+sortType)
	}

	return &SortParam{
		Sort: strings.Join(param, ","),
	}
}

func Sorting(ctx ServiceContextManager, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var (
		sort string = ctx.GetSortParam().Sort
	)

	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sort)
	}
}
