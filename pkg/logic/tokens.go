package logic

import (
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"shop/pkg/constants"
	"shop/pkg/models"
	"time"
)

func CreateAccessToken(email string) (string, error){
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&models.AuthTokenClaim{
			&jwt.StandardClaims{
				ExpiresAt: time.Now().Add(constants.AccessTokenExpireTime).Unix(),
			},
			models.Email{email},
		},
	)

	return accessToken.SignedString([]byte(constants.SigningToken))
}

func CreateRefreshToken() string {
	b := make([]byte, constants.RefreshTokenLength)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}