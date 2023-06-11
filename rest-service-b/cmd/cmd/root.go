package cmd

import (
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/api"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/config"

	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "rest-service-b",
	Short: "Another simple rest api to test some concepts",
	Long:  `Another simple rest api to test some concepts.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		config.Setup(cfgFile)
		logger.SetupLogs(config.GetLogLevel(), config.GetLogFormat(), config.GetServiceName())
		telemetry.InitTelemetry(
			telemetry.WithEnvironment(config.GetEnvironment()),
			telemetry.WithMetricsEndpoint(config.GetOtelMetricsEndpoint()),
			telemetry.WithTracesEndpoint(config.GetOtelTraceEndpoint()),
			telemetry.WithVersion(config.GetVersion()),
			telemetry.WithServiceName(config.GetServiceName()),
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.Start(apiPort)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	apiPort int
	cfgFile string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rest-api.yaml)")

	rootCmd.PersistentFlags().IntVarP(&apiPort, "port", "p", 8080, "Port to API listen for requests")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
