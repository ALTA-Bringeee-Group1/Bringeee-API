package web

type Pagination struct {
	Page       		int `json:"page"`
	Limit      		int `json:"limit"`
	TotalPages 		int `json:"total_pages"`
	TotalRecords	int `json:"total_records"`
}
