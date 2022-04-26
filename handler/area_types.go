package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
	"github.com/ONSdigital/log.go/v2/log"
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

	if err := req.Valid(); err != nil {
		h.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "invalid query parameters"),
		)
		return
	}

	cReq := cantabular.GetGeographyDimensionsRequest{
		Dataset: req.Dataset,
	}
	cReq.Limit = req.Limit
	cReq.Offset = req.Offset

	res, err := h.ctblr.GetGeographyDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			errors.Wrap(err, "failed to get area-types"),
		)
		return
	}

	log.Info(ctx, "Got response from Cantabular", log.Data{
		"response": res,
	})

	var resp contract.GetAreaTypesResponse

	if res != nil {
		resp.PaginationResponse = contract.PaginationResponse{
			Count:      res.Count,
			TotalCount: res.TotalCount,
		}
		resp.Limit = res.Limit
		resp.Offset = res.Offset
		for _, edge := range res.Dataset.RuleBase.IsSourceOf.Edges {
			resp.AreaTypes = append(resp.AreaTypes, model.AreaType{
				ID:         edge.Node.Name,
				Label:      edge.Node.Label,
				TotalCount: edge.Node.Categories.TotalCount,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}
