package exceptions

import "errors"

type Warning struct {
	Message string
}

func (e *Warning) Error() string {
	return e.Message
}

func NewWarning(message string) *Warning {
	return &Warning{Message: message}
}

func IsWarning(err error) bool {
	var service *Warning
	return errors.As(err, &service)
}