package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

type PaginationParams struct {
	Limit  int `json:"limit" schema:"limit"`
	Offset int `json:"offset" schema:"offset"`
}
type PaginationResponse struct {
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
	PaginationParams
}

// GetAreaTypesRequest defines the schema for the GET /area-types query parameter
type GetAreaTypesRequest struct {
	Dataset string `schema:"dataset"`
	PaginationParams
}

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	AreaTypes []model.AreaType `json:"area-types"`
	PaginationResponse
}
