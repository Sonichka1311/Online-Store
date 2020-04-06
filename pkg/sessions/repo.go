package sessions

import (
	"net/http"
	"shop/pkg/models"
)

type Repo struct {
	Connector *models.Connector
}

var (
	SessionHandler = "session"
)

func (r *Repo) GetUserFrom(data *[]byte) (*http.Response, *models.Error) {
	return r.Connector.Post(SessionHandler, data)
}

func (r *Repo) AddOrUpdate(session *Session) (*http.Response, *models.Error) {
	data, jsonError := session.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Put(SessionHandler, data)
}

func (r *Repo) Delete(session *Session) (*http.Response, *models.Error) {
	data, jsonError := session.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Delete(SessionHandler, data)
}
