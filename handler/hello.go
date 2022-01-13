package handler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"errors"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
)

type Hello struct{
	respond Responder
}

func NewHello(r Responder) *Hello {
	return &Hello{
		respond: r,
	}
}

// Get is the handler for GET /hello
func(h *Hello) Get(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	resp := contract.GetHelloResponse{
		Message: "Hello, World!",
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// Post is the handler for POST /hello - Is used for an error example
func(h *Hello) Post(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	defer r.Body.Close()

	var req contract.PostHelloRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond.Error(ctx, w, Error{
			err:        fmt.Errorf("badly formed request body: %w", err),
			statusCode: http.StatusBadRequest,
		})
		return
	}

	if req.Error {
		h.respond.Error(ctx, w, Error{
			err:        errors.New("I am logged error"),
			resp:       "Hello, error!",
			statusCode: http.StatusUnauthorized,
			logData:    log.Data{
				"hello": "world",
			},
		})
		return
	}

	h.respond.StatusCode(w, http.StatusOK)
}
