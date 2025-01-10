package common

type QueryParam struct {
	PerPage string `query:"per_page" validate:"omitempty,numeric"`
	Page    string `query:"page" validate:"omitempty,numeric"`
	Sort    string `query:"sort"`
}

type ParamID struct {
	ID int `param:"id" validate:"required,numeric"`
}

type ParamCode struct {
	Code string `param:"code" validate:"required"`
}
