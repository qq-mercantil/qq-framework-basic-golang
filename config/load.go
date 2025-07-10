package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/caarlos0/env/v6"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
)

type Option func(*env.Options)

// These constants help create conditions to specific environments
const (
	Production  = "production"
	Staging     = "staging"
	Development = "development"
	Test        = "test"
)

func LoadGeneric[T any](appName string, reference *T) error {
    sLog := logger.Get()

    if len(appName) > 0 {
        os.Setenv("APP_NAME", appName)
    }

    if err := env.Parse(reference); err != nil {
        sLog.Errorf("error parsing configs - %+v", err)
        return err
    }

    appEnv := os.Getenv("APP_ENV")
    if !isValidEnvironment(appEnv) {
        return fmt.Errorf("invalid environment: %s", appEnv)
    }

    return nil
}

func LoadMultiple(appName string, configs ...any) error {
    sLog := logger.Get()

    if len(appName) > 0 {
        os.Setenv("APP_NAME", appName)
    }

    for _, config := range configs {
        if err := env.Parse(config); err != nil {
            sLog.Errorf("error parsing config %T - %+v", config, err)
            return err
        }
    }

    appEnv := os.Getenv("APP_ENV")
    if !isValidEnvironment(appEnv) {
        return fmt.Errorf("invalid environment: %s", appEnv)
    }

    return nil
}

func isValidEnvironment(env string) bool {
	valid := slices.Contains([]string{Production, Staging, Development, Test}, env)

	return valid
}
