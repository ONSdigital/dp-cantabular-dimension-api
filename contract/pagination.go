package contract

const (
	defaultLimit = 20
)

type PaginationParams struct {
	Limit  int `json:"limit" schema:"limit"`
	Offset int `json:"offset" schema:"offset"`
}
type PaginationResponse struct {
	PaginationParams
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}
