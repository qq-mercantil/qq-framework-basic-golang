package application

import (
	"context"

	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
)

type abstractUseCase[P any, E exceptions.Exception, R any] interface {
	OnExecute(ctx context.Context, props P) (R, E)
}

type BaseUseCase[P any, E exceptions.Exception, R any] struct {
	processName string
	executor    abstractUseCase[P, E, R]
    *ApplicationProcess[P, E, R]
}

func NewBaseUseCase[P any, E exceptions.Exception, R any](
    processName string,
    executor abstractUseCase[P, E, R],
    applicationProcessLoggerService ApplicationProcessLoggerService,
) *BaseUseCase[P, E, R] {
    baseUseCase := &BaseUseCase[P, E, R]{
        processName: processName,
        executor:    executor,
    }

    applicationProcess := NewApplicationProcess(processName, executor, applicationProcessLoggerService)
    baseUseCase.ApplicationProcess = applicationProcess

    return baseUseCase
}
