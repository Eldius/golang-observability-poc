package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/httpclient"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/config"
	"io"
	"net/url"
)

func GetWeather(ctx context.Context, city string) (*Weather, error) {
	ctx, closer := telemetry.StartSpan(ctx, "GetWeatherIntegration")
	defer closer()

	telemetry.AddAttribute(ctx, "city", city)

	l := logger.GetLogger(ctx)

	endpoint, err := url.Parse(config.GetWeatherServiceEndpoint())
	if err != nil {
		telemetry.NotifyError(ctx, err)
		l.Error().Err(err).Msg("error parsing endpoint")
		return nil, err
	}
	q := endpoint.Query()
	q.Set("city", city)

	endpoint.RawQuery = q.Encode()

	l.Debug().Str("api_key", config.GetWeatherServiceApiKey()).Msg("integrating")
	resp, err := httpclient.GetRequest(ctx, endpoint.String(), httpclient.WithHeader("x-api-key", config.GetWeatherServiceApiKey()))
	if err != nil {
		telemetry.NotifyError(ctx, err)
		l.Error().Err(err).Msg("error requesting weather api")
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("integration error %d: %s", resp.StatusCode, b)
	}
	var w Weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		telemetry.NotifyError(ctx, err)
		return nil, err
	}

	return &w, nil
}
