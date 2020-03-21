package repository

import (
	"bytes"
	"net/http"
)

type Connector struct {
	Host string
}

func (connector Connector) Get(handler string) (*http.Response, error) {
	return http.Get(connector.Host + handler)
}

func (connector Connector) Post(handler string, data []byte) (*http.Response, error) {
	return http.Post(connector.Host + handler, "application/json", bytes.NewBuffer(data))
}

