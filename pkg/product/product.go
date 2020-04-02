package product

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

type AllItems struct {
	Items       []Product `json:"items"`
	PagesCount  int       `json:"pages_count"`
	CurrentPage int       `json:"page"`
}
