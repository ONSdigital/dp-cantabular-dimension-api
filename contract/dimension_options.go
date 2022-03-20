package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetDimensionOptionsRequest defines the schema for the GET /dimension/{dimension}/options query parameter
type GetDimensionOptionsRequest struct {
	Dimension string `json:"dimension"`
}

// GetDimensionOptionsResponse is the response object for GET /dimension/{dimension}/options
type GetDimensionOptionsResponse struct {
	Options []*model.DimensionOption `json:"options"`
}
