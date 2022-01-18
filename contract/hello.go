package contract

import (
	"github.com/ONSdigital/dp-cantabular-dimension-api/model"
)

type PostHelloRequest struct{
	model.Hello
}

type GetHelloResponse struct{
	Message string `json:"message"`
}
