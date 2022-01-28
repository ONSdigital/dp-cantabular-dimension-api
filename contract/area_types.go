package contract

import "github.com/ONSdigital/dp-cantabular-dimension-api/model"

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	AreaTypes []model.AreaType `json:"area-types"`
}
