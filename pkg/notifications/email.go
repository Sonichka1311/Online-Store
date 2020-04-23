package notifications

import (
	"log"
	"net/smtp"
	"shop/pkg/constants"
	"strconv"
	"strings"
)

func SendEmail(to string, msg string) bool {
	auth := smtp.PlainAuth("", constants.TestUser, constants.TestPassword, constants.MockServer)

	err := smtp.SendMail(
		strings.Join([]string{constants.MockServer, strconv.Itoa(constants.MockPort)}, ":"),
		auth,
		constants.TestUser,
		[]string{to},
		[]byte(msg))

	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return false
	}
	return true
}
