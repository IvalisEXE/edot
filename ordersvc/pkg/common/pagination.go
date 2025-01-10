package common

type Pagination struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}
