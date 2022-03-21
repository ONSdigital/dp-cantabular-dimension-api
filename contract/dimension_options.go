package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetDimensionOptionsResponse is the response object for GET /dimension/{dimension}/options
type GetDimensionOptionsResponse struct {
	Options []model.DimensionOption `json:"options"`
}
