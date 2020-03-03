package models

type Product struct {
    Name      string `json:"name"`
    Id        int `json:"id"`
    Category  string `json:"category"`
}