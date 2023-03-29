package telemetry

import (
    "context"
    "github.com/rs/zerolog/log"
)

type options struct {
    serviceName     string
    version         string
    env             string
    metricsEndpoint string
    tracesEndpoint  string
    ctx             context.Context
}

func (o *options) MetricsEnabled() bool {
    return o.metricsEndpoint != ""
}

func (o *options) TracesEnabled() bool {
    return o.tracesEndpoint != ""
}

type Option func(o *options) *options

func InitTelemetry(opts ...Option) {
    opt := &options{
        serviceName:     "default-name",
        version:         "v0.0.0",
        env:             "dev",
        metricsEndpoint: "",
        tracesEndpoint:  "",
        ctx:             context.Background(),
    }
    for _, f := range opts {
        opt = f(opt)
    }
    log.Debug().Msg("starting telemetry providers")

    initMetrics(opt)
    initTracer(opt)

    log.Debug().Msg("ending telemetry providers")
}

// WithServiceName configures service name
func WithServiceName(n string) Option {
    return func(o *options) *options {
        o.serviceName = n
        return o
    }
}

func WithVersion(v string) Option {
    return func(o *options) *options {
        o.version = v
        return o
    }
}

func WithEnvironment(e string) Option {
    return func(o *options) *options {
        o.env = e
        return o
    }
}

func WithMetricsEndpoint(e string) Option {
    return func(o *options) *options {
        o.metricsEndpoint = e
        return o
    }
}

func WithTracesEndpoint(e string) Option {
    return func(o *options) *options {
        o.tracesEndpoint = e
        return o
    }
}

func WithContext(c context.Context) Option {
    return func(o *options) *options {
        o.ctx = c
        return o
    }
}
