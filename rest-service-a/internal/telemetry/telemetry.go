package telemetry

import (
	"context"
	"github.com/rs/zerolog/log"
)

func InitTelemetry() {
	log.Debug().Msg("starting telemetry providers")

	ctx := context.Background()

	initMetrics(ctx)
	initTracer(ctx)

	log.Debug().Msg("ending telemetry providers")
}
