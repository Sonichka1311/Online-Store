package models

import (
	"encoding/json"
	"log"
	"net/http"
)

type Replier struct {
	Writer *http.ResponseWriter
}

func (replier *Replier) ReplyWithMessage(message string) error {
	err := json.NewEncoder(*(replier.Writer)).Encode(
		ReplyWithMessage{
			Message: message,
		},
	)
	if err != nil {
		log.Printf("Reply with message err: %s\n", err.Error())
		return err
	}
	return nil
}

func (replier *Replier) ReplyWithJson(data interface{}) error {
	return json.NewEncoder(*(replier.Writer)).Encode(data)
}

func (replier *Replier) ReplyWithData(data []byte) error {
	_, writeError := (*replier.Writer).Write(data)
	return writeError
}

func (replier *Replier) ReplyWithError(err *Error) {
	// ToDo: log
	(*replier.Writer).WriteHeader(err.ErrorCode)
	writeError := json.NewEncoder(*replier.Writer).Encode(err)
	if writeError != nil {
		// ToDo: log
		http.Error(*replier.Writer, writeError.Error(), http.StatusInternalServerError)
	}
}