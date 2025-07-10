package application

import (
	"encoding/json"
	"fmt"

	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
	"github.com/qq-mercantil/qq-framework-basic-golang/utils"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
)

type ApplicationProcessOffLoggerService struct {
	logger *logger.Event
}

func NewApplicationProcessOffLoggerService() *ApplicationProcessOffLoggerService {
	sLog := logger.Get()

	sLog.Debug("ApplicationProcessOffLoggerService constructor")

	return &ApplicationProcessOffLoggerService{
		logger: sLog,
	}
}

func (s *ApplicationProcessOffLoggerService) Start(processName string, props string) (string, *exceptions.ApplicationProcessException) {
	formattedProps := utils.FormatPrintableMessage(props, 1000)

	s.logger.Debugf("starting %s with parameters %s", processName, formattedProps)

	return "", nil
}

func (s *ApplicationProcessOffLoggerService) Success(result any, processName, applicationProcessLogId string) (*struct{}, *exceptions.ApplicationProcessException) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		resultJSON = []byte(fmt.Sprintf("%+v", result))
	}

	formattedResult := utils.FormatPrintableMessage(string(resultJSON), 1000)

	s.logger.Debugf("process %s executed with SUCCESS result %s", processName, formattedResult)

	return nil, nil
}

func (s *ApplicationProcessOffLoggerService) Error(err exceptions.Exception, processName string, applicationProcessLogId string) (*struct{}, *exceptions.ApplicationProcessException) {
	s.logger.Errorf("%s executed with ERROR with result %v", processName, err)

	return nil, nil
}
