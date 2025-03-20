package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       int 	`json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}