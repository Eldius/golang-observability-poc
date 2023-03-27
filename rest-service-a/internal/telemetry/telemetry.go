package telemetry

import (
	"context"
	"github.com/rs/zerolog/log"
)

func InitTelemetry() {
	log.Debug().Msg("starting tracer provider")

	ctx := context.Background()

	initTracer(ctx)
}
