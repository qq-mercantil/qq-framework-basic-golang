package application

import "github.com/qq-mercantil/qq-framework-basic-golang/exceptions"


type ApplicationProcessLoggerService interface {
	Start(processName, props string) (string, *exceptions.ApplicationProcessException)
	Error(err exceptions.Exception, processName, applicationProcessId string) (*struct{}, *exceptions.ApplicationProcessException)
	Success(result any, processName, applicationProcessId string) (*struct{}, *exceptions.ApplicationProcessException)
}
