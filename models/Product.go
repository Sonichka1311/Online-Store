package models

type Product struct {
    Name      string `json:"name"`
    Id        int `json:"id"`
    Category  string `json:"category"`
}

func (product *Product) SetName(name string) {
    product.Name = name
}

func (product *Product) SetId(id int) {
    product.Id = id
}

func (product *Product) SetCategory(category string) {
    product.Category = category
}