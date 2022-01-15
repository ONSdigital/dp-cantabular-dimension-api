package contract

type PostHelloRequest struct{
	CantabularBlob string `json:"blob"`
}

type GetHelloResponse struct{
	Message string `json:"message"`
}
