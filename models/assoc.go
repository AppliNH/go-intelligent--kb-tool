package models

import (
	"encoding/json"
)

// Association describes an association of name and surname
type Association struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Pk      float64 `json:"pk"`
}

// JSONTOAssociation creates a Association object from a JSON string
func JSONTOAssociation(jsonStr string) Association {
	var a Association
	json.Unmarshal([]byte(jsonStr), &a)

	return a
}
