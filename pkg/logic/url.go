package logic

import (
	"strconv"
	"strings"
)

func GetUrl(protocol, host string, port int) string {
	return strings.Join([]string{protocol, "://", host, ":", strconv.Itoa(port)}, "")
}
