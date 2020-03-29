package commonModels

import (
	"encoding/json"
	"net/http"
)

type Replier struct {
	Writer *http.ResponseWriter
}

func (replier *Replier) ReplyWithMessage(message string) error {
	return json.NewEncoder(*(replier.Writer)).Encode(
		ReplyWithMessage{
			Message: message,
		},
	)
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