// Package dto defines the data transfer objects exchanged with the Vue
// frontend. GORM models are never serialized directly — every page gets a
// DTO describing exactly the shape of its props, and every write endpoint
// gets a DTO describing (and validating) its input.
package dto

// PaginationDTO describes the paging state of a list page.
type PaginationDTO struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"perPage"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}
