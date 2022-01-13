package responder

// ErrorResponse is the generic ONS error response for HTTP errors
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
