package application

import (
	"context"
	"encoding/json"

	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
)

type APPLICATION_PROCESS_LOGGER_SERVICE_TYPE string

var (
	POSTGRES APPLICATION_PROCESS_LOGGER_SERVICE_TYPE = "POSTGRES"
	OFF      APPLICATION_PROCESS_LOGGER_SERVICE_TYPE = "OFF"
)

type ApplicationProcess[P any, E exceptions.Exception, R any] struct {
	processName                     string
	applicationProcessLoggerService ApplicationProcessLoggerService
	logger                          *logger.Event
	useCase                         abstractUseCase[P, E, R]
}

func NewApplicationProcess[P any, E exceptions.Exception, R any](
    processName string,
    useCase abstractUseCase[P, E, R],
    applicationProcessLoggerService ApplicationProcessLoggerService, 
) *ApplicationProcess[P, E, R] {
    return &ApplicationProcess[P, E, R]{
        processName:                     processName,
        applicationProcessLoggerService: applicationProcessLoggerService,
        logger:                          logger.Get(),
        useCase:                         useCase,
    }
}

func (p *ApplicationProcess[P, E, R]) Execute(ctx context.Context, props P) (R, exceptions.Exception) {
	var zero R

	propsString, err := json.Marshal(props)
	if err != nil {
		p.logger.Error("Erro ao serializar propriedades:", err)

		return zero, exceptions.NewApplicationProcessException("Failed to marshal properties")
	}

	processLogId, logErr := p.applicationProcessLoggerService.Start(p.processName, string(propsString))
	if logErr != nil {
		p.logger.Error("Erro ao iniciar log:", logErr)
		return zero, exceptions.NewApplicationProcessException("Failed to start log")
	}

	result, processErr := p.useCase.OnExecute(ctx, props)

	p.generateFinalLog(processLogId, processErr, result)

	return result, processErr
}

func (p *ApplicationProcess[P, E, R]) generateFinalLog(processLogId string, processErr exceptions.Exception, result R) {
	if processErr != nil {
		p.applicationProcessLoggerService.Error(processErr, p.processName, processLogId)
	} else {
		p.applicationProcessLoggerService.Success(result, p.processName, processLogId)
	}
}
