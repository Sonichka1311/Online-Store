package product

import (
	"fmt"
	"net/http"
	"shop/pkg/models"
)

type Repo struct {
	Connector *models.Connector
}

var (
	GetAll = func() string { return "" }
	Get = func(id string) string { return fmt.Sprintf("product/%s", id) }
	Add =  func() string { return "product" }
	Edit =  func() string { return "product" }
	Delete = func() string { return "product" }
)

func (r *Repo) Get(id string) (*http.Response, *models.Error) {
	return r.Connector.Get(Get(id))
}

func (r *Repo) GetAll() (*http.Response, *models.Error) {
	return r.Connector.Get(GetAll())
}

func (r *Repo) Add(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Post(Add(), data)
}

func (r *Repo) Edit(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Put(Edit(), data)
}

func (r *Repo) Delete(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Delete(Delete(), data)
}
