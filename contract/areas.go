package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetAreasRequest defines the schema for the GET /areas query parameter

type GetAreasRequest struct {
	Dataset        string
	DimensionNames []string
}

// GetAreasResponse is the response object for GET /areas
type GetAreasResponse struct {
	Areas []model.Areas `json:"areas"`
}
