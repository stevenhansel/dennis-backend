package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	LISTEN_ADDR             string `mapstructure:"LISTEN_ADDR"`
	POSTGRES_CONNECTION_URI string `mapstructure:"POSTGRES_CONNECTION_URI"`
}

func initializeDevelopmentConfig() error {
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	return nil
}

func initializeProductionConfig(config *Configuration) error {
	fields := reflect.VisibleFields(reflect.TypeOf(*config))
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		for _, field := range fields {
			fieldName := field.Tag.Get("mapstructure")
			if fieldName == pair[0] {
				reflect.ValueOf(config).Elem().FieldByName(fieldName).SetString(pair[1])
			}
		}
	}

	return nil
}

func New(environment Environment) (*Configuration, error) {
	var config Configuration

	if environment == DEVELOPMENT {
		if err := initializeDevelopmentConfig(); err != nil {
			return nil, err
		}
	} else {
		if err := initializeProductionConfig(&config); err != nil {
			return nil, err
		}
	}

	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
