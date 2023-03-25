package config

import "github.com/spf13/viper"

func GetServiceName() string {
	return "rest-api"
}

func GetLogLevel() string {
	return viper.GetString("log.level")
}

func GetOtelEndpoint() string {
	return viper.GetString("otel.endpoint")
}
