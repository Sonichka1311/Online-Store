package constants

import "time"

const (
	AccessTokenExpireTime = time.Minute * 5
	RefreshTokenExpireTime = time.Hour
	SigningToken = "2816e66cb08c9f4cb5d7c080b2fca85f17cdb1cbe32380c7fdde9cf469185e30"
	RefreshTokenLength = 20
	Protocol = "http"
	MainHost = "localhost"
	DatabaseHost = "tarantool"
	AuthHost = "auth"
	MainServerPort = 8080
	DatabasePort = 3301
	AuthPort = 3687
)
