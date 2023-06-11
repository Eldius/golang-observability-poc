package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/httpclient"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/rest-service-b/internal/config"
	"github.com/pkg/errors"
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
		err = errors.Wrap(err, "failed to parse endpoint")
		telemetry.NotifyError(ctx, err)
		l.WithError(err).Error("error parsing endpoint")
		return nil, err
	}
	q := endpoint.Query()
	q.Set("city", city)

	endpoint.RawQuery = q.Encode()

	l.WithField("api_key", config.GetWeatherServiceAPIKey()).Info("integrating")
	resp, err := httpclient.GetRequest(ctx, endpoint.String(), httpclient.WithHeader("x-api-key", config.GetWeatherServiceAPIKey()))
	if err != nil {
		err = errors.Wrap(err, "failed to call external weather integration")
		telemetry.NotifyError(ctx, err)
		l.WithError(err).Error("error requesting weather api")
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body) //nolint:errcheck // ignoring error
		return nil, errors.New(fmt.Sprintf("integration error %d: %s", resp.StatusCode, b))
	}
	var w Weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		err = errors.Wrap(err, "failed to decode external weather integration response")
		telemetry.NotifyError(ctx, err)
		return nil, err
	}

	return &w, nil
}
