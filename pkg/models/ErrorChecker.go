package models

import "log"

type Error struct {
	ErrorCode 	int 	`json:"code"`
	ErrorString string	`json:"message"`
}

func NewError(err error, code int) (*Error, bool) {
	if err != nil {
		return &Error{
			ErrorCode: code,
			ErrorString: err.Error(),
		}, true
	}
	return nil, false
}

type ErrorChecker struct {
	Error 	*Error
	Replier *Replier
}

func (checker *ErrorChecker) ReplyError() {
	log.Printf("Replying user with error: %s\n", checker.Error.ErrorString)
	(*checker.Replier).ReplyWithError(checker.Error)
	checker.Error = nil
}

func (checker *ErrorChecker) CheckCustomError(err error, code int) bool {
	if newError, ok := NewError(err, code); ok {
		checker.Error = newError
		checker.ReplyError()
		return true
	}
	return false
}

func (checker *ErrorChecker) CheckError(err *Error) bool {
	if err != nil {
		checker.Error = err
		checker.ReplyError()
		return true
	}
	return false
}

func (checker *ErrorChecker) NewError(message string, code int) bool {
	checker.Error = &Error{
		ErrorString: message,
		ErrorCode:   code,
	}
	checker.ReplyError()
	return true
}
