package models

// Association describes an association of name and surname
type Association struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Pk      float64 `json:"pk"`
}
