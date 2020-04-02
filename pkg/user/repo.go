package user

import (
	"net/http"
	"shop/pkg/models"
)

type Repo struct {
	Connector *models.Connector
}

var (
	Add = func() string { return "user" }
	Get = func() string { return "session" }
	Update = func() string { return "session" }
)

func (r *Repo) Add(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Post(Add(), data)
}

func (r *Repo) Get(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Post(Get(), data)
}

func (r *Repo) Update(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Put(Update(), data)
}
