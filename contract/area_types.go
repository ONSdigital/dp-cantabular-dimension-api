package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetAreaTypesRequest defines the schema for the GET /area-types query parameter
type GetAreaTypesRequest struct {
	CantabularDataset string `schema:"cantabular_dataset"`
}

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	AreaTypes []model.AreaType `json:"area-types"`
}
