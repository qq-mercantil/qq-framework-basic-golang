package exceptions

import "errors"

type InvalidPropsException struct {
	Message string
}

func (e *InvalidPropsException) Error() string {
	return e.Message
}

func NewInvalidPropsException(message string) *InvalidPropsException {
	return &InvalidPropsException{Message: message}
}

func IsInvalidPropsException(err error) bool {
	var invalidProps *InvalidPropsException
	return errors.As(err, &invalidProps)
}
