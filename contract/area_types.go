package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetAreaTypesRequest defines the schema for the GET /area-types query parameter
type GetAreaTypesRequest struct {
	Dataset string `schema:"dataset"`
}

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	AreaTypes []model.AreaType `json:"area-types"`
}

// GetAreaTypeRequest defines the schema for the GET /area-types/{area-type} query parameter
type GetAreaTypeRequest struct {
	Dataset  string `schema:"dataset"`
	AreaType string `schema:"area-type"`
}

// GetAreaTypeResponse is the response object for GET /area-types/{area-type}
type GetAreaTypeResponse struct {
	AreaType model.AreaType `json:"area-type"`
}
