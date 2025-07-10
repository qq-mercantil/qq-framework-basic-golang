package exceptions

import "errors"

type RepositoryException struct {
	message string
}

func NewRepositoryException(message string) RepositoryException {
	return RepositoryException{
		message: message,
	}
}

func (e RepositoryException) Error() string {
	return e.message
}

func IsRepositoryException(err error) bool {
	var invalidProps *RepositoryException
	return errors.As(err, &invalidProps)
}
