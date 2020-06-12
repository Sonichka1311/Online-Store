package constants

import "time"

var (
	AccessTokenExpireTime       = time.Minute * 5
	RefreshTokenExpireTime      = time.Hour
	ConfirmationTokenExpireTime = time.Minute * 10
	SigningToken                = "2816e66cb08c9f4cb5d7c080b2fca85f17cdb1cbe32380c7fdde9cf469185e30"
	MainHost                    = "localhost"
	Protocol                    = "http"
	NotificationUrl             = "http://a5f2055f8d20.ngrok.io"
	AuthHost                    = "auth"
	MainPort                    = 8080
	AuthPort                    = 3687
	ValidatePort                = 3867
	UploadPort                  = 3838
	QueueConnectionRetries      = 10
	QueueConnectionSleepTime    = time.Second * 10
	QueueNotificationServer     = "amqp://guest:guest@rabbitmq:5672/"
	QueueUploadServer           = "amqp://guest:guest@rabbitmqupload:5762/"
	MockServer                  = "mock"
	TestUser                    = "test@mock"
	TestPassword                = "1234"
	MockPort 					= 25
	SuperAdmin			 		= "admin"
	DatabaseConnectionRetries 	= 10
	DatabaseConnectionSleepTime = time.Second * 10
)
