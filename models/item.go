package models

// Item represents an item in the bolt DB
type Item struct {
	Type  string  `json:"type"`
	Value string  `json:"value"`
	Pk    float64 `json:"pk"`
}
