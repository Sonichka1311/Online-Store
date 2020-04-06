package user

import (
	"fmt"
	"net/http"
	"shop/pkg/auth"
	"shop/pkg/models"
)

type Repo struct {
	Connector *models.Connector
}

var (
	UserHandler = "user"
	ConfirmationHandler = "confirmation"
	GetConfirmationHandler = func(token string) string { return fmt.Sprintf("%s/%s", ConfirmationHandler, token) }
	GetUserHandler = func(email string) string { return fmt.Sprintf("%s/%s", UserHandler, email)}
)

func (r *Repo) AddUser(usr *User) (*http.Response, *models.Error) {
	data, jsonError := usr.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Post(UserHandler, data)
}

func (r *Repo) GetUser(usr *User) (*http.Response, *models.Error) {
	return r.Connector.Get(GetUserHandler(usr.Email))
}

func (r *Repo) AddConfirmation(confirmation *auth.Confirmation) (*http.Response, *models.Error) {
	data, jsonError := confirmation.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Post(ConfirmationHandler, data)
}

func (r *Repo) GetConfirmation(token string) (*http.Response, *models.Error) {
	return r.Connector.Get(GetConfirmationHandler(token))
}

func (r *Repo) DeleteConfirmation(confirmation *auth.Confirmation) (*http.Response, *models.Error) {
	data, jsonError := confirmation.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Delete(ConfirmationHandler, data)
}

func (r *Repo) ConfirmUser(usr *User) (*http.Response, *models.Error) {
	data, jsonError := usr.GetJson()
	if jsonError != nil {
		return nil, jsonError
	}
	return r.Connector.Put(UserHandler, data)
}
