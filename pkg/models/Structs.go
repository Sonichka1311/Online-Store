package models

import (
	"github.com/dgrijalva/jwt-go"
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

type EmailNotification struct {
	Email 	string
	Message string
}
