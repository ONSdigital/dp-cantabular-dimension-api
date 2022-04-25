package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
}

type cantabularClient interface {
	GetCodebook(ctx context.Context, req cantabular.GetCodebookRequest) (*cantabular.GetCodebookResponse, error)
	GetGeographyDimensions(ctx context.Context, dataset string) (*cantabular.GetGeographyDimensionsResponse, error)
	GetAreas(ctx context.Context, req cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)
}

type validator interface {
	Valid() error
}
