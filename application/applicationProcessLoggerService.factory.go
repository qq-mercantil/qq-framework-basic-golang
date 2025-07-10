package application

import (
	"sync"

	"gorm.io/gorm"
)

var (
	pgLoggerInstance  *ApplicationProcessPGLoggerService
	offLoggerInstance *ApplicationProcessOffLoggerService
	pgLoggerOnce      sync.Once
	offLoggerOnce     sync.Once
)

func FactoryCreate(serviceType APPLICATION_PROCESS_LOGGER_SERVICE_TYPE, db *gorm.DB) ApplicationProcessLoggerService {
	if serviceType == POSTGRES {
		pgLoggerOnce.Do(func() {
			pgLoggerInstance = NewApplicationProcessPGLoggerService(db)
		})
		return pgLoggerInstance
	}

	offLoggerOnce.Do(func() {
		offLoggerInstance = NewApplicationProcessOffLoggerService()
	})
	return offLoggerInstance
}

func ResetLoggerInstances() {
	pgLoggerInstance = nil
	offLoggerInstance = nil
	pgLoggerOnce = sync.Once{}
	offLoggerOnce = sync.Once{}
}
