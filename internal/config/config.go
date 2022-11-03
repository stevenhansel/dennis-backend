package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Parse everything as a string initially
type RawConfiguration struct {
	LISTEN_ADDR             string `mapstructure:"LISTEN_ADDR"`
	POSTGRES_CONNECTION_URI string `mapstructure:"POSTGRES_CONNECTION_URI"`
	CORS_ORIGINS            string `mapstructure:"CORS_ORIGINS"`
}

type Configuration struct {
	LISTEN_ADDR             string
	POSTGRES_CONNECTION_URI string
	CORS_ORIGINS            []string
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

func initializeProductionConfig(config *RawConfiguration) error {
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
	var rawConfig RawConfiguration
	var config Configuration

	if environment == DEVELOPMENT {
		if err := initializeDevelopmentConfig(); err != nil {
			return nil, err
		}

		err := viper.Unmarshal(&rawConfig)
		if err != nil {
			return nil, err
		}

	} else {
		if err := initializeProductionConfig(&rawConfig); err != nil {
			return nil, err
		}
	}

	rawConfigMap := map[string]string{}
	rawVal := reflect.ValueOf(&rawConfig).Elem()
	for i := 0; i < rawVal.NumField(); i++ {
		fieldName := rawVal.Type().Field(i).Name
		rawConfigMap[fieldName] = rawVal.FieldByName(fieldName).String()
	}

	cfgVal := reflect.ValueOf(&config).Elem()
	for i := 0; i < cfgVal.NumField(); i++ {
		fieldName := cfgVal.Type().Field(i).Name

		var value string
		if v, ok := rawConfigMap[fieldName]; ok {
			value = v
		}

		if value == "" {
			continue
		}

		t := cfgVal.Field(i).Interface()
		switch t.(type) {
		case []string:
			slice := strings.Split(strings.Trim(value, " "), ",")
			for j := 0; j < len(slice); j++ {
				slice[j] = strings.Trim(slice[j], " ")
			}

			cfgVal.Field(i).Set(reflect.ValueOf(slice))
			break
		case string:
			cfgVal.Field(i).SetString(value)
			break
		}
	}

	return &config, nil
}
