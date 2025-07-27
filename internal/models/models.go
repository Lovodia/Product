package models

import "time"

type Category struct {
	ID   int    `json:"id"`
	Name string `json: "name"`
}

type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"category_at"`
}
