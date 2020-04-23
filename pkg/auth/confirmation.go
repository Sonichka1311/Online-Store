package auth

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/models"
	"time"
)

type Confirmation struct {
	Email 	string 	`json:"email"`
	Token 	string 	`json:"token"`
	Expire 	int64 	`json:"expire"`
}

func (c *Confirmation) Init(email, token string) {
	c.Email = email
	c.Token = token
	c.Expire = time.Now().Add(constants.ConfirmationTokenExpireTime).Unix()
}

func (c *Confirmation) GetConfirmation(body io.ReadCloser) *models.Error {
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if err, ok := models.NewError(bodyParseError, http.StatusBadRequest); ok {
		return err
	}

	parsedBodyError := json.Unmarshal(readBody, &c)
	if err, ok := models.NewError(parsedBodyError, http.StatusBadRequest); ok {
		return err
	}

	return nil
}

func (c *Confirmation) GetJson() (*[]byte, *models.Error) {
	jsonData, jsonError := json.Marshal(c)
	if err, ok := models.NewError(jsonError, http.StatusBadRequest); ok {
		return nil, err
	}
	return &jsonData, nil
}


