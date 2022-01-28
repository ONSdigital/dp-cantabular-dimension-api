package contract

import (
	"errors"

	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
)

// CreateHelloRequest is the request object for POST /hello
type CreateHelloRequest struct{
	model.Hello
}

// Valid satisfies the Validator interface which allows incoming requests
// to be be made to parse a basic validation check
func (r *CreateHelloRequest) Valid() error {
	if len(r.CantabularBlob) < 1{
		return errors.New("missing/empty field: 'cantabular_blob'")
	}
	return nil
}

// GetHelloResponse is the response object for GET /hello
type GetHelloResponse struct{
	Message string `json:"message"`
}
