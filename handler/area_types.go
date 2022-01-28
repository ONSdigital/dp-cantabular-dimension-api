package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
)

// AreaTypes handles requests to /area-types
type AreaTypes struct {
	respond responder
	ctblr   cantabularClient
}

// NewAreaTypes returns a new AreaTypes handler
func NewAreaTypes(r responder, c cantabularClient) *AreaTypes {
	return &AreaTypes{
		respond: r,
		ctblr:   c,
	}
}

// Get is the handler for GET /area-types
func (at *AreaTypes) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := contract.GetAreaTypesResponse{
		AreaTypes: []model.AreaType{
			{
				ID:    "id1",
				Label: "label1",
			},
		},
	}

	at.respond.JSON(ctx, w, http.StatusOK, resp)
}
