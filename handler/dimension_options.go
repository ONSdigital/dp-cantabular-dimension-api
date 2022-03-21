package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// DimensionOptions handles requests to dimension options
type DimensionOptions struct {
	respond responder
	ctblr   cantabularClient
}

// NewDimensionOptions returns a new DimensionOptions handler
func NewDimensionOptions(r responder, c cantabularClient) *DimensionOptions {
	return &DimensionOptions{
		respond: r,
		ctblr:   c,
	}
}

// Get is the handler for GET /dimension/{dimension}/options
func (s *DimensionOptions) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dimension := chi.URLParam(r, "dimension")
	cantabularRequest := cantabular.GetDimensionOptionsRequest{
		DimensionNames: []string{dimension},
	}
	cantabularResponse, err := s.ctblr.GetDimensionOptions(ctx, cantabularRequest)
	if err != nil {
		s.respond.Error(
			ctx,
			w,
			dperrors.StatusCode(err),
			errors.Wrap(err, "failed to get dimension options"),
		)
		return
	}

	var response contract.GetDimensionOptionsResponse
	if cantabularResponse != nil {
		for _, dimension := range cantabularResponse.Dataset.Table.Dimensions {
			for _, category := range dimension.Categories {
				option := model.DimensionOption{
					Name: dimension.Variable.Name,
					Links: model.DimensionOptionLinks{
						Code: model.Link{
							ID: category.Code,
						},
					},
				}
				response.Options = append(response.Options, option)
			}
		}
	}

	s.respond.JSON(ctx, w, http.StatusOK, response)
}
