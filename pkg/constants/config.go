package constants

import "time"

var (
	AccessTokenExpireTime       = time.Minute * 20
	RefreshTokenExpireTime      = time.Hour * 24
	ConfirmationTokenExpireTime = time.Minute * 10
	SigningToken                = "2816e66cb08c9f4cb5d7c080b2fca85f17cdb1cbe32380c7fdde9cf469185e30"
	MainHost                    = "server"
	Protocol                    = "http"
	NotificationUrl             = "a5f2055f8d20.ngrok.io"
	AuthHost                    = "auth"
	MainPort                    = 8080
	AuthPort                    = 3687
	ValidatePort                = 3867
	UploadPort                  = 3838
	QueueConnectionRetries      = 10
	QueueConnectionSleepTime    = time.Second * 10
	QueueNotificationServer     = "amqp://guest:guest@rabbitmq:5672/"
	QueueUploadServer           = "amqp://guest:guest@rabbitmqupload:5672/"
	MockServer                  = "mock"
	TestUser                    = "test@mock"
	TestPassword                = "1234"
	MockPort 					= 25
	SuperAdmin			 		= "admin"
	DatabaseConnectionRetries 	= 10
	DatabaseConnectionSleepTime = time.Second * 10
	SmsRuId 					= "F043B94F-3F15-BA86-6720-FA6537CB03B6"
)
