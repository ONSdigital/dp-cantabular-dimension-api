package handler

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
)

// Hello handles requests to /hello
type Hello struct{
	respond responder
	ctblr   cantabularClient
}

// NewHello returns a new Hello handler
func NewHello(r responder, c cantabularClient) *Hello {
	return &Hello{
		respond: r,
		ctblr:   c,
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
func(h *Hello) Create(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var req contract.CreateHelloRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond.Error(ctx, w, Error{
			err:        fmt.Errorf("badly formed request body: %w", err),
			message:    "badly formed request body",
			statusCode: http.StatusBadRequest,
		})
		return
	}
	defer r.Body.Close()

	cReq := cantabular.GetCodebookRequest{
		DatasetName: req.CantabularBlob,
		Variables:   []string{"sex", "city", "siblings_3"},
		Categories:  false,
	}

	resp, err := h.ctblr.GetCodebook(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, fmt.Errorf("failed to get Codebook: %w", err))
		return
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}
