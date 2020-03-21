package logic

import (
	"../models"
	"../repository"
	"net/http"
	"sync"
)

var mutex = sync.Mutex{}
var connector = repository.Connector{Host: "http://172.18.0.2:8888"}

func Get(handler string)(*http.Response, *models.Error) {
	mutex.Lock()
	defer mutex.Unlock()
	return Request(handler, "get", nil)
}

func Post(handler string, data *[]byte)(*http.Response, *models.Error) {
	mutex.Lock()
	defer mutex.Unlock()
	return Request(handler, "post", data)
}

func Request(handler string, requestType string, data *[]byte)(*http.Response, *models.Error) {
	var response *http.Response
	var requestError error
	switch requestType {
	case "get":
		response, requestError = connector.Get("/" + handler)
	case "post":
		response, requestError = connector.Post("/" + handler, *data)
	}
	if requestError != nil {
		return response, &models.Error{
			ErrorString: requestError.Error(),
			ErrorCode: http.StatusInternalServerError}
	}
	if response.StatusCode != http.StatusOK {
		return response, &models.Error{
			ErrorString: response.Status,
			ErrorCode: response.StatusCode}
	}
	return response, nil
}

