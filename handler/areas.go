package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
	"github.com/gorilla/schema"

	"github.com/pkg/errors"
)

// Areas handles requests to /area-types
type Areas struct {
	respond responder
	ctblr   cantabularClient
}

// NewAreas returns a new Areas handler
func NewAreas(r responder, c cantabularClient) *Areas {
	return &Areas{
		respond: r,
		ctblr:   c,
	}
}

// Get is the handler for GET /areas
func (h *Areas) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req cantabular.QueryData
	if err := schema.NewDecoder().Decode(&req, r.URL.Query()); err != nil {
		h.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to decode query parameters"),
		)
		return
	}

	res, err := h.ctblr.GetAreas(ctx, req)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			errors.Wrap(err, "failed to get areas"),
		)
		return
	}

	var resp contract.GetAreasResponse

	if res != nil {
		for _, edge := range res.Dataset.RuleBase.IsSourceOf.CategorySearch.Edges {
			resp.Areas = append(resp.Areas, model.Areas{
				ID:       edge.Node.Code,
				Label:    edge.Node.Label,
				AreaType: edge.Node.Variable.Name,
			})

		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}
