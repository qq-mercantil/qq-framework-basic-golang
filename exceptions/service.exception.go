package exceptions

import "errors"

type ServiceException struct {
	Message string
}

func (e *ServiceException) Error() string {
	return e.Message
}

func NewServiceException(message string) *ServiceException {
	return &ServiceException{Message: message}
}

func IsServiceException(err error) bool {
	var service *ServiceException
	return errors.As(err, &service)
}