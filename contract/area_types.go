package contract

import (
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"

	"github.com/pkg/errors"
)

// GetAreaTypesRequest defines the schema for the GET /area-types query parameter
type GetAreaTypesRequest struct {
	PaginationParams
	Dataset string `schema:"dataset"`
}

// Valid validates the values given in the request
func (r *GetAreaTypesRequest) Valid() error {
	if r.Limit < 0 {
		return errors.New("'limit' must be greater than 0")
	}

	if r.Offset < 0 {
		return errors.New("'offset' must be greater than 0")
	}

	if r.Limit == 0 {
		r.Limit = defaultLimit
	}

	return nil
}

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	PaginationResponse
	AreaTypes []model.AreaType `json:"area-types"`
}
