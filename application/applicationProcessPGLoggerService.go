package application

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
	"github.com/qq-mercantil/qq-framework-basic-golang/model"
	"github.com/qq-mercantil/qq-framework-basic-golang/utils"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
	"gorm.io/gorm"
)

type ApplicationProcessPGLoggerService struct {
	logger *logger.Event
	db     *gorm.DB
}

func NewApplicationProcessPGLoggerService(db *gorm.DB) *ApplicationProcessPGLoggerService {
	sLog := logger.Get()

	sLog.Debug("ApplicationProcessPGLoggerService constructor")

	return &ApplicationProcessPGLoggerService{
		logger: sLog,
		db:     db,
	}
}

func (s *ApplicationProcessPGLoggerService) Start(processName string, props string) (string, *exceptions.ApplicationProcessException) {
	formattedProps := utils.FormatPrintableMessage(props, 1000)

	s.logger.Debugf("starting %s with parameters %s", processName, formattedProps)

	id := uuid.New().String()
	processLog := &model.ProcessLogModel{
		ID:          id,
		ProcessName: processName,
		Props:       &formattedProps,
		MessageType: model.BEGINNING,
		StartedAt:   time.Now(),
	}

	if err := processLog.Save(); err != nil {
		s.logger.Errorf("Failed to save process log for %s: %v", processName, err)
		return "", exceptions.NewApplicationProcessException(fmt.Sprintf("Failed to save process log", err))
	}

	return "", nil
}

func (s *ApplicationProcessPGLoggerService) Success(result any, processName, applicationProcessLogId string) (*struct{}, *exceptions.ApplicationProcessException) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		resultJSON = []byte(fmt.Sprintf("%+v", result))
	}

	formattedResult := utils.FormatPrintableMessage(string(resultJSON), 1000)

	s.logger.Debugf("process %s executed with SUCCESS result %s", processName, formattedResult)

	now := time.Now()
	updates := map[string]interface{}{
		"message_type": model.FINISH_WITH_SUCCESS,
		"result":       formattedResult,
		"finished_at":  now,
	}

	if err := model.UpdateByID(s.db, applicationProcessLogId, updates); err != nil {
		s.logger.Errorf("Failed to update process log %s: %v", applicationProcessLogId, err)
		return nil, exceptions.NewApplicationProcessException("Failed to update process log")
	}

	return nil, nil
}

func (s *ApplicationProcessPGLoggerService) Error(err exceptions.Exception, processName string, applicationProcessLogId string) (*struct{}, *exceptions.ApplicationProcessException) {
	var errorStack string

	if stackTracer, ok := err.(interface{ StackTrace() string }); ok {
		errorStack = stackTracer.StackTrace()
	} else if goErr, ok := err.(error); ok {
		errorStack = goErr.Error()
	} else {
		errorStack = fmt.Sprintf("%+v", err)
	}

	s.logger.Errorf("%s executed with ERROR with result %s", processName, errorStack)

	now := time.Now()
	errorClass := model.GetErrorClass(err)
	errorType := model.GetErrorType(err)

	updates := map[string]interface{}{
		"message_type": model.FINISH_WITH_ERROR,
		"error":        utils.FormatPrintableMessage(errorStack, 1000),
		"error_class":  errorClass,
		"error_type":   errorType,
		"finished_at":  now,
	}

	if updateErr := model.UpdateByID(s.db, applicationProcessLogId, updates); updateErr != nil {
		s.logger.Errorf("Failed to update process log %s: %v", applicationProcessLogId, updateErr)
		return nil, exceptions.NewApplicationProcessException("Failed to update process log")
	}

	return nil, nil
}
