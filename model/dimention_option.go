package model

// LinkObject represents a generic structure for all links
type LinkObject struct {
	HRef string `json:"href"`
	ID   string `json:"id"`
}

// DimensionOption contains unique information and metadata used when processing the data
type DimensionOption struct {
	Name  string               `json:"name"`
	Links DimensionOptionLinks `json:"links"`
}

// DimensionOptionLinks represents a list of link objects related to dimension options
type DimensionOptionLinks struct {
	Code     LinkObject `json:"code"`
	CodeList LinkObject `json:"code_list"`
	Version  LinkObject `json:"version"`
}
