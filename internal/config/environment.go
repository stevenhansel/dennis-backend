package config

import "fmt"

type Environment string

const (
	DEVELOPMENT Environment = "development"
	STAGING     Environment = "staging"
	PRODUCTION  Environment = "production"
)

var availableEnvironments = map[string]Environment{
	"development": DEVELOPMENT,
	"staging":     STAGING,
	"production":  PRODUCTION,
}

func (e *Environment) Set(value string) error {
	env, ok := availableEnvironments[value]
	if !ok {
		return fmt.Errorf("value not valid")
	}

	*e = env

	return nil
}

func (e *Environment) String() string {
	return string(*e)
}
