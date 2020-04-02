package models

import (
	"bytes"
	"errors"
	"net/http"
)

type Router struct {
	Host string
}

func (router Router) Get(handler string) (*http.Response, error) {
	return http.Get(router.Host + handler)
}

func (router Router) Post(handler string, data []byte) (*http.Response, error) {
	return http.Post(router.Host + handler, "application/json", bytes.NewBuffer(data))
}

func (router Router) Put(handler string, data []byte) (*http.Response, error) {
	return router.PutAndDelete(handler, data, http.MethodPut)
}

func (router Router) Delete(handler string, data []byte) (*http.Response, error) {
	return router.PutAndDelete(handler, data, http.MethodDelete)
}

func (router Router) PutAndDelete(handler string, data []byte, method string) (*http.Response, error) {
	client := &http.Client{}
	request, requestError := http.NewRequest(method, router.Host + handler, bytes.NewBuffer(data))
	if requestError != nil {
		return nil, errors.New("500 Internal server error")
	}
	response, execError := client.Do(request)
	if execError != nil {
		return nil, errors.New("500 Internal server error")
	}
	return response, nil
}