package model

// Link represents a generic structure for all links
type Link struct {
	HRef string `json:"href"`
	ID   string `json:"id"`
}

// DimensionOption contains unique information and metadata used when processing the data
type DimensionOption struct {
	Name  string               `json:"name"`
	Links DimensionOptionLinks `json:"links"`
}

// DimensionOptionLinks represents a list of links related to dimension options
type DimensionOptionLinks struct {
	Code     Link `json:"code"`
	CodeList Link `json:"code_list"`
	Version  Link `json:"version"`
}
