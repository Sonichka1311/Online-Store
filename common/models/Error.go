package commonModels

type ErrorChecker struct {
	Error 	*Error
	Replier *Replier
}

func (checker *ErrorChecker) ReplyError() {
	(*checker.Replier).ReplyWithError(checker.Error)
	checker.Error = nil
}

func (checker *ErrorChecker) CheckCustomError(err error, code int) bool {
	if err != nil {
		checker.Error = &Error{ErrorString: err.Error(), ErrorCode: code}
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
