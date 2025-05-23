package dto

type PaginationRequest struct {
	Limit  int    `json:"limit" query:"limit" validate:"required,gt=0"`
	Page   int    `json:"page" query:"page" validate:"required,gt=0"`
	SortBy string `json:"sort_by" query:"sort_by,omitempty"`
}

type PaginationResponse struct {
	BaseResponse
	Pagination *PaginationMetadata `json:"pagination"`
	Link       *Link               `json:"link"`
}

type Link struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type PaginationMetadata struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
}
