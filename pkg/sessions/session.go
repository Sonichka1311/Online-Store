package sessions

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/models"
	"time"
)

type Session struct {
	Email 		 string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	Expire       int64  `json:"expire"`
}

func (s *Session) Init(email, token string) {
	s.Email = email
	s.RefreshToken = token
	s.Expire = time.Now().Add(constants.RefreshTokenExpireTime).Unix()
}

func (s *Session) Update(token string) {
	s.RefreshToken = token
	s.Expire = time.Now().Add(constants.RefreshTokenExpireTime).Unix()
}

func (s *Session) GetSession(body io.ReadCloser) *models.Error {
	defer body.Close()
	readBody, bodyParseError := ioutil.ReadAll(body)
	if err, ok := models.NewError(bodyParseError, http.StatusBadRequest); ok {
		return err
	}

	parsedBodyError := json.Unmarshal(readBody, &s)
	if err, ok := models.NewError(parsedBodyError, http.StatusBadRequest); ok {
		return err
	}

	return nil
}

func (s *Session) GetJson() (*[]byte, *models.Error) {
	jsonData, jsonError := json.Marshal(s)
	if err, ok := models.NewError(jsonError, http.StatusBadRequest); ok {
		return nil, err
	}
	return &jsonData, nil
}
