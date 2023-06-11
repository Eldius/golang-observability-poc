package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

var (
	logger *logrus.Entry
)

func GetLogger(ctx context.Context) *logrus.Entry {
	span := trace.SpanFromContext(ctx)

	return logger.WithFields(logrus.Fields{
		"trace_id": span.SpanContext().TraceID().String(),
		"span_id":  span.SpanContext().SpanID().String(),
	})
}

func Logger() *logrus.Entry {
	return logger
}

func SetupLogs(logLevel, logFormat, service string) {
	var logFormatter logrus.Formatter
	// Log as JSON instead of the default ASCII formatter.
	if strings.ToLower(logFormat) == "json" {
		logFormatter = &logrus.JSONFormatter{}
	} else {
		logFormatter = &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		}
	}

	logrus.SetFormatter(logFormatter)
	logrus.SetReportCaller(true)

	logLevel = strings.ToLower(logLevel)
	switch strings.ToLower(logLevel) {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

	hostname, _ := os.Hostname()
	var standardFields = logrus.Fields{
		"hostname": hostname,
		"service":  service,
	}
	logger = logrus.StandardLogger().WithFields(standardFields)

	logger.WithField("setup_log_level", logrus.GetLevel()).Info("SetupLogsEnd")
}
