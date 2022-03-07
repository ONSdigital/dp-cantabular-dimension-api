package model

// Areas is an area model with ID, Label, and area-type
type Areas struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	AreaType string `json:"area-type"`
}
