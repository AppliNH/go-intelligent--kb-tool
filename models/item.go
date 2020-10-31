package models

import "encoding/json"

// Item represents an item in the bolt DB
type Item struct {
	Type  string  `json:"type"`
	Value string  `json:"value"`
	Pk    float64 `json:"pk"`
}

// JSONTOItem creates a Item object from a JSON string
func JSONTOItem(jsonStr string) Item {
	var a Item
	json.Unmarshal([]byte(jsonStr), &a)

	return a
}
