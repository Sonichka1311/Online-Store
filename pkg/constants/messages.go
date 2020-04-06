package constants

import "fmt"

var (
	ConfirmationMessage = func(token string) string {
		return fmt.Sprintf("To verify your account, click http://%s:%d/%s\n", MainHost, AuthPort, token)
	}
	SignUpOkMessage = func(email string) string {
		return fmt.Sprintf("Confirmation request was sent to %s\n", email)
	}
	ConfirmOkMessage = "Account verified"
	ValidAccessToken = "Access token is valid"
)