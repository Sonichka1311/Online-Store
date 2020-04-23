package product

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"shop/pkg/models"
)

type Product struct {
	Name      string `json:"name"`
	Id        int `json:"id"`
	Category  string `json:"category"`
}

func (product Product) GetJson() ([]byte, *models.Error) {
	jsonData, jsonError := json.Marshal(product)
	if err, isErr := models.NewError(jsonError, http.StatusInternalServerError); isErr {
		return nil, err
	}
	return jsonData, nil
}

func (product *Product) GetFromBody(reader io.ReadCloser) *models.Error {
	defer reader.Close()
	body, bodyParseError := ioutil.ReadAll(reader)
	if err, isErr := models.NewError(bodyParseError, http.StatusInternalServerError); isErr {
		return err
	}
	unmarshalError := json.Unmarshal(body, product)
	if err, isErr := models.NewError(unmarshalError, http.StatusInternalServerError); isErr {
		return err
	}

	return nil
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
	Items       []*Product `json:"items"`
	PagesCount  int        `json:"pages_count"`
	CurrentPage int        `json:"page"`
}
