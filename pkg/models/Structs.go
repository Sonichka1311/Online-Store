package models

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"net/http"
)

type ReplyWithMessage struct {
	Message string `json:"message"`
}

type AccessToken struct {
	TokenStr string `json:"access_token"`
}

type ReplyWithTokens struct {
	AccessToken  string
	RefreshToken string
}

type Email struct {
	Email string `json:"email"`
}

type AuthTokenClaim struct {
	*jwt.StandardClaims
	Email
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Notification struct {
	Email 	string
	Phone   string
	Message string
}

type Verification struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

func (v Verification) GetJson() ([]byte, *Error) {
	jsonData, jsonError := json.Marshal(v)
	if err, isErr := NewError(jsonError, http.StatusInternalServerError); isErr {
		return nil, err
	}
	return jsonData, nil
}

func (v *Verification) GetFromBody(reader io.ReadCloser) *Error {
	defer reader.Close()
	body, bodyParseError := ioutil.ReadAll(reader)
	if err, isErr := NewError(bodyParseError, http.StatusInternalServerError); isErr {
		return err
	}
	unmarshalError := json.Unmarshal(body, v)
	if err, isErr := NewError(unmarshalError, http.StatusInternalServerError); isErr {
		return err
	}

	return nil
}
