package handler

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/ONSdigital/dp-cantabular-dimension-api/contract"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

type Hello struct{
	respond Responder
	ctblr   CantabularClient
}

func NewHello(r Responder, c CantabularClient) *Hello {
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

	h.respond.Raw(ctx, w, http.StatusOK, []byte(resp.Message))
}

// Post is the handler for POST /hello - Is used for an error example
func(h *Hello) Post(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var req contract.PostHelloRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond.Error(ctx, w, Error{
			err:        fmt.Errorf("badly formed request body: %w", err),
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
