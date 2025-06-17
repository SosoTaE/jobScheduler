package structs

import "jobScheduler/models"

type PaginationMeta struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
}

type PaginatedJobsResponse struct {
	Data []models.Job   `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
