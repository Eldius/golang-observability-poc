package config

import "github.com/spf13/viper"

var (
    Version string
)

func GetServiceName() string {
    return "rest-service-b"
}

func GetLogLevel() string {
    return viper.GetString("log.level")
}

func GetDBHost() string {
    return viper.GetString("db.host")
}

func GetDBPort() int {
    return viper.GetInt("db.port")
}

func GetDBUser() string {
    return viper.GetString("db.user")
}

func GetDBPass() string {
    return viper.GetString("db.pass")
}

func GetDBName() string {
    return viper.GetString("db.database")
}

func GetDBSSLEnabled() bool {
    return viper.GetBool("db.ssl")
}

func GetMigrationsEnabled() bool {
    return viper.GetBool("db.migrations.enabled")
}

func GetMigrationsSource() string {
    return viper.GetString("db.migrations.source")
}

func GetEnvironment() string {
    return viper.GetString("environment")
}

func GetVersion() string {
    if Version == "" {
        return "v0.0.0"
    }

    return Version
}

func GetWeatherServiceEndpoint() string {
    return viper.GetString("integration.serviceb.endpoint")
}
