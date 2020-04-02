package models

type ReplyWithMessage struct {
	Message string `json:"message"`
}

type Error struct {
	ErrorCode 	int 	`json:"code"`
	ErrorString string	`json:"message"`
}

type AccessToken struct {
	TokenStr string `json:"access_token"`
}
