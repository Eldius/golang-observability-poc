package config

import "github.com/spf13/viper"

func GetOtelTraceEndpoint() string {
	return viper.GetString("otel.trace.endpoint")
}

func GetOtelMetricsEndpoint() string {
	return viper.GetString("otel.metrics.endpoint")
}

func EnableTraceRest() bool {
	return viper.GetBool("telemetry.rest.enable")
}

func EnableTraceDB() bool {
	return viper.GetBool("telemetry.db.enable")
}

func OTLPUser() string {
	return viper.GetString("telemetry.auth.user")
}

func OTLPPass() string {
	return viper.GetString("telemetry.auth.pass")
}
