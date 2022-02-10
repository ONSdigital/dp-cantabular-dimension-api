package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"

	dperrors "github.com/ONSdigital/dp-net/v2/errors"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
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
func (h *AreaTypes) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req contract.GetAreaTypesRequest
	if err := schema.NewDecoder().Decode(&req, r.URL.Query()); err != nil {
		h.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to decode query parameters"),
		)
		return
	}

	res, err := h.ctblr.GetGeographyDimensions(ctx, req.Dataset)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			dperrors.StatusCode(err), // Can be changed to ctblr.StatusCode(err) once added to Client
			errors.Wrap(err, "failed to get area-types"),
		)
		return
	}

	var resp contract.GetAreaTypesResponse

	if res != nil {
		for _, edge := range res.Dataset.RuleBase.IsSourceOf.Edges {
			resp.AreaTypes = append(resp.AreaTypes, model.AreaType{
				ID:    edge.Node.Name,
				Label: edge.Node.Label,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}
