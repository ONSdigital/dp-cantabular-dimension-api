package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

type coder interface {
	Code() int
}

// Responder handles responding to http requests
type Responder interface{
	JSON(context.Context,http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, error) 
	StatusCode(http.ResponseWriter, int)
}

type CantabularClient interface {
	GetCodebook(context.Context, cantabular.GetCodebookRequest) (*cantabular.GetCodebookResponse, error)
}
