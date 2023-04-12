package config

import (
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// Setup initializes app configuration.
func Setup(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".docker-profiler-go" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("rest-api-config")
	}
	setDefaults()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("api")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		l := logger.Logger()
		l.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func setDefaults() {
	viper.SetDefault("integration.weather.endpoint", "https://api.api-ninjas.com/v1/weather")
	viper.SetDefault("log.level", "DEBUG")
}
