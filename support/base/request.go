package base

type PaginationRequest struct {
	Sort     string `json:"sort" form:"sort"`
	Includes string `json:"includes" form:"includes"`
	Page     int    `json:"page" form:"page" binding:"omitempty,min=1"`
	PerPage  int    `json:"per_page" form:"per_page" binding:"omitempty,min=1"`
}
