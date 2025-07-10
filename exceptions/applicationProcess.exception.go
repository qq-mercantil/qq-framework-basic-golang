package exceptions

import "errors"

type ApplicationProcessException struct {
	Message string
}

func (e *ApplicationProcessException) Error() string {
	return e.Message
}

func NewApplicationProcessException(message string) *ApplicationProcessException {
	return &ApplicationProcessException{Message: message}
}

func IsApplicationProcessException(err error) bool {
	var applicationProcessException *ApplicationProcessException
	return errors.As(err, &applicationProcessException)
}