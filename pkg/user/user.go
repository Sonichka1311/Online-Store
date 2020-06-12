package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"shop/pkg/models"
)

type User struct {
	Email 		string 	`json:"email"`
	Password 	string 	`json:"password,omitempty"`
	Confirm		bool	`json:"confirm,omitempty"`
	Role		string  `json:"role,omitempty"`
}

func (u *User) GetUser(body io.ReadCloser) *models.Error {
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if err, ok := models.NewError(bodyParseError, http.StatusBadRequest); ok {
		log.Println("Fail to read body")
		return err
	}

	parsedBodyError := json.Unmarshal(readBody, &u)
	if err, ok := models.NewError(parsedBodyError, http.StatusBadRequest); ok {
		log.Println("Fail to parse body into user")
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
