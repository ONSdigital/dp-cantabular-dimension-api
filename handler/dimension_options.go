package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	"github.com/gorilla/schema"
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

	var request contract.GetDimensionOptionsRequest
	if err := schema.NewDecoder().Decode(&request, r.URL.Query()); err != nil {
		s.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to decode query parameters"),
		)
		return
	}

	cantabularRequest := cantabular.GetDimensionOptionsRequest{
		DimensionNames: []string{request.Dimension},
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
						Code: model.LinkObject{
							ID: category.Code,
						},
					},
				}
				response.Options = append(response.Options, &option)
			}
		}
	}

	s.respond.JSON(ctx, w, http.StatusOK, response)
}
