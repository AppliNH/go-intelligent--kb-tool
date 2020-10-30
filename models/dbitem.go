package models

// DBItem represents an item in the bolt DB
type DBItem struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
