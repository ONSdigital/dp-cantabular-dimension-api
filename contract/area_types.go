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

// GetAreaTypeParentsRequest defines the schema for the GET /area-types/{area-type}/parents query parameter
type GetAreaTypeParentsRequest struct {
	Dataset  string `schema:"dataset"`
	AreaType string `schema:"area-type"`
}

// GetAreaTypeParentsResponse is the response object for GET /area-types/{area-type}/parents
type GetAreaTypeParentsResponse struct {
	AreaTypes []model.AreaType `json:"area-types"`
}
