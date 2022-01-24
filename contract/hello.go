package contract

import (
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
)

// CreateHelloRequest is the request object for POST /hello
type CreateHelloRequest struct{
	model.Hello
}

// GetHelloResponse is the response object for GET /hello
type GetHelloResponse struct{
	Message string `json:"message"`
}
