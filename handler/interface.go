package handler

import (
	"context"
	"net/http"
)

// Responder handles responding to http requests
type Responder interface{
	JSON(context.Context,http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, error) 
	StatusCode(http.ResponseWriter, int)
}
