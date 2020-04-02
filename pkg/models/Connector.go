package models

import (
	"net/http"
	"sync"
)

type Connector struct {
	Router Router
	Mutex  sync.Mutex
}

func (connector Connector) Get(handler string) (*http.Response, *Error) {
	connector.Mutex.Lock()
	defer connector.Mutex.Unlock()
	return connector.Request(handler, http.MethodGet, nil)
}

func (connector Connector) Post(handler string, data *[]byte) (*http.Response, *Error) {
	connector.Mutex.Lock()
	defer connector.Mutex.Unlock()
	return connector.Request(handler, http.MethodPost, data)
}

func (connector Connector) Put(handler string, data *[]byte) (*http.Response, *Error) {
	connector.Mutex.Lock()
	defer connector.Mutex.Unlock()
	return connector.Request(handler, http.MethodPut, data)
}

func (connector Connector) Delete(handler string, data *[]byte) (*http.Response, *Error) {
	connector.Mutex.Lock()
	defer connector.Mutex.Unlock()
	return connector.Request(handler, http.MethodDelete, data)
}

func (connector Connector) Request(handler string, requestType string, data *[]byte) (*http.Response, *Error) {
	var response *http.Response
	var requestError error
	switch requestType {
	case http.MethodGet:
		response, requestError = connector.Router.Get("/" + handler)
	case http.MethodPost:
		response, requestError = connector.Router.Post("/" + handler, *data)
	case http.MethodPut:
		response, requestError = connector.Router.Put("/" + handler, *data)
	case http.MethodDelete:
		response, requestError = connector.Router.Delete("/" + handler, *data)
	}
	if requestError != nil {
		return response, &Error{
			ErrorString: requestError.Error(),
			ErrorCode: http.StatusInternalServerError}
	}
	if response.StatusCode != http.StatusOK {
		return response, &Error{
			ErrorString: response.Status,
			ErrorCode: response.StatusCode}
	}
	return response, nil
}

