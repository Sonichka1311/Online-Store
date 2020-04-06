package logic

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"shop/pkg/constants"
	"shop/pkg/models"
	"time"
)

func CreateAccessToken(email string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": models.Email{Email: email},
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	return accessToken.SignedString([]byte(constants.SigningToken))
}

func CreateRefreshToken() (*string, error) {
	uuidToken, uuidError := uuid.NewRandom()
	if uuidError != nil {
		return nil, uuidError
	}
	token := uuidToken.String()

	return &token, nil
}

func CreateConfirmationToken() (*string, error) {
	uuidToken, uuidError := uuid.NewRandom()
	if uuidError != nil {
		return nil, uuidError
	}
	token := uuidToken.String()

	return &token, nil
}
