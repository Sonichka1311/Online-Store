package models

import "github.com/dgrijalva/jwt-go"

type ReplyWithTokens struct {
	AccessToken  string
	RefreshToken string
}

type Email struct {
	Email string `json:"email"`
}

type User struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type Session struct {
	Email 		 string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	Expire       string `json:"expire"`
}

type AuthTokenClaim struct {
	*jwt.StandardClaims
	Email
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	TokenStr string `json:"refresh_token"`
}

type RefreshTokenUpdate struct {
	OldToken string `json:"old_refresh_token"`
	Token 	 string `json:"refresh_token"`
	Expire   string `json:"expire"`
	Time     string `json:"time"`
}
