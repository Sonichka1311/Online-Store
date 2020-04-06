package constants

import "time"

const (
	AccessTokenExpireTime = time.Minute * 5
	RefreshTokenExpireTime = time.Hour
	ConfirmationTokenExpireTime = time.Minute * 10
	SigningToken = "2816e66cb08c9f4cb5d7c080b2fca85f17cdb1cbe32380c7fdde9cf469185e30"
	Protocol = "http"
	MainHost = "localhost"
	DatabaseHost = "tarantool"
	AuthHost = "auth"
	MainServerPort = 8080
	DatabasePort = 3301
	AuthPort = 3687
	QueueConnectionRetries = 10
	QueueConnectionSleepTime = time.Second * 5
	QueueServer = "amqp://guest:guest@rabbitmq:5672/"
	SmtpServer = "" // TODO
	TestUser = "test@mock"
	TestPassword = "1234"
	MockPort = 3031
)
