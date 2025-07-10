package application

import (
	"github.com/qq-mercantil/qq-framework-db-golang/db"
	"go.uber.org/fx"
)

func Module(loggerType APPLICATION_PROCESS_LOGGER_SERVICE_TYPE, customImpl *ApplicationProcessLoggerService) fx.Option {
	return fx.Module(
        "application",
        fx.Provide(fx.Annotate(
            func(db *db.GormDatabase) ApplicationProcessLoggerService {
				if customImpl != nil {
					return *customImpl
				}

                return FactoryCreate(loggerType, db.DB)
            }, 
            fx.As(new(ApplicationProcessLoggerService)),
        )),
    )
}
