package exceptions

import "errors"

type BusinessException struct {
	Message string
}

func (e *BusinessException) Error() string {
	return e.Message
}

func NewBusinessException(message string) *BusinessException {
	return &BusinessException{Message: message}
}

func IsBusinessException(err error) bool {
	var businessException *BusinessException
	return errors.As(err, &businessException)
}
