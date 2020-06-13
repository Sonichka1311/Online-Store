package product

import (
	"encoding/json"
	"encoding/xml"
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
	Name      string `json:"name",xml:"Name"`
	Id        int `json:"id",xml:"Id"`
	Category  string `json:"category",xml:"Category"`
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
	//log.Printf("LINE: %s\n", line)

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
	//log.Printf("Category: %s\n", split[2])
	product.Category = strings.Trim(split[2], " \r\r\t\n")
	return nil
}

func (product *Product) GetFromXML(lines string) *models.Error {
	// convert lines to full xml
	lines = "<?xml version=\"1.0\"?>\n<Tag>\n" + lines + "</Tag>"

	// structure for parsing
	type XML struct {
		XMLName   xml.Name `xml:"Tag"`
		Name      string `xml:"Name"`
		Id        int `xml:"Id"`
		Category  string `xml:"Category"`
	}

	// parse xml
	var x XML
	err := xml.Unmarshal([]byte(lines), &x)
	if err != nil {
		log.Println(err.Error())
		newErr, _ :=  models.NewError(err, http.StatusBadRequest)
		return newErr
	}

	// init product
	product.Name = x.Name
	product.Category = x.Category

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
