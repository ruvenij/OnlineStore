package model

type PaginationParams struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
