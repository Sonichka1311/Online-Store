package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"shop/pkg/models"
)

type User struct {
	Email 		string 	`json:"email"`
	Password 	string 	`json:"password,omitempty"`
	Confirm		bool	`json:"confirm,omitempty"`
}

func (u *User) GetUser(body io.ReadCloser) *models.Error {
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if err, ok := models.NewError(bodyParseError, http.StatusBadRequest); ok {
		return err
	}

	parsedBodyError := json.Unmarshal(readBody, &u)
	if err, ok := models.NewError(parsedBodyError, http.StatusBadRequest); ok {
		return err
	}

	return nil
}

func (u *User) GetJson() (*[]byte, *models.Error) {
	jsonData, jsonError := json.Marshal(u)
	if err, ok := models.NewError(jsonError, http.StatusBadRequest); ok {
		return nil, err
	}
	return &jsonData, nil
}
