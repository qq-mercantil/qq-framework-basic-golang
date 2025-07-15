package exceptions

import "errors"

type RepositoryNoDataFoundException struct {
	message string
}

func NewRepositoryNoDataFoundException(message string) *RepositoryNoDataFoundException {
	return &RepositoryNoDataFoundException{
		message: message,
	}
}

func (e RepositoryNoDataFoundException) Error() string {
	return e.message
}

func IsRepositoryNoDataFoundException(err error) bool {
	var repositoryNoDataFound *RepositoryNoDataFoundException
	return errors.As(err, &repositoryNoDataFound)
}
