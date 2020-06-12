package product

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"shop/pkg/models"
	"strconv"
	"strings"
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

func (product *Product) GetFromCsv(line string) *models.Error {
	log.Printf("LINE: %s\n", line)

	// split line by columns
	split := strings.Split(line, ";")
	if len(split) != 3 {
		if err, isErr := models.NewError(errors.New("Bad line: "+line), http.StatusInternalServerError); isErr {
			return err
		}
	}

	// try to create product by split values
	ans, parseErr := strconv.ParseUint(strings.Trim(split[0], " \t\n"), 10, 64)
	if parseErr != nil {
		if err, isErr := models.NewError(parseErr, http.StatusInternalServerError); isErr {
			return err
		}
	}
	product.Id = int(ans)
	product.Name = strings.Trim(split[1], " \t\n")
	product.Category = strings.Trim(split[2], " \t\n")
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

type Uploading struct {
	Token 		string
	Products 	[]Product
}
