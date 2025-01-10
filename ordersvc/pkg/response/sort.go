package response

import (
	"strings"

	"gorm.io/gorm"

	"ordersvc/pkg/common"
)

func AssertSortPayload(p common.QueryParam) *common.SortParam {
	var param []string

	if p.Sort == "" {
		return &common.SortParam{}
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

	return &common.SortParam{
		Sort: strings.Join(param, ","),
	}
}

func Sorting(ctx common.ServiceContextManager, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var (
		sort string = ctx.GetSortParam().Sort
	)

	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sort)
	}
}
