package model

// Areas is an area model with ID and Label
type Areas struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	AreaType string `json:"area-type"`
}
