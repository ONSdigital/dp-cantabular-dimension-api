package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

type coder interface {
	Code() int
}

type responder interface{
	JSON(context.Context,http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, error) 
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
}

type cantabularClient interface {
	GetCodebook(context.Context, cantabular.GetCodebookRequest) (*cantabular.GetCodebookResponse, error)
}
