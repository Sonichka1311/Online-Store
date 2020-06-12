package logic

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"shop/pkg/constants"
	"shop/pkg/models"
	"shop/pkg/user"
	"time"
)

var ReturnedAccessToken string // ONLY for tests
func CreateAccessToken(usr *user.User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": models.Email{Email: usr.Email},
		"role": usr.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(constants.AccessTokenExpireTime).Unix(),
	})

	ReturnedAccessToken , _ = accessToken.SignedString([]byte(constants.SigningToken))
	return accessToken.SignedString([]byte(constants.SigningToken))
}

var ReturnedRefreshToken string // ONLY for tests
func CreateRefreshToken() (*string, error) {
	uuidToken, uuidError := uuid.NewRandom()
	if uuidError != nil {
		log.Println("Fail to generate refresh token")
		return nil, uuidError
	}
	token := uuidToken.String()

	ReturnedRefreshToken = token
	return &token, nil
}

var ReturnedConfirmationToken string // ONLY for tests
func CreateConfirmationToken() (*string, error) {
	uuidToken, uuidError := uuid.NewRandom()
	if uuidError != nil {
		log.Println("Fail to generate confirmation token")
		return nil, uuidError
	}
	token := uuidToken.String()

	ReturnedConfirmationToken = token
	return &token, nil
}
