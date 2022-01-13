package contract

type PostHelloRequest struct{
	Error bool `json:"error"`
}

type GetHelloResponse struct{
	Message string `json:"message"`
}
