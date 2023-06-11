package serviceb

import (
	"context"
	"encoding/json"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/httpclient"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/rest-service-a/internal/config"
	"github.com/pkg/errors"
	"net/url"
)

func GetWeather(ctx context.Context, city string) (*Weather, error) {

	e, err := url.Parse(config.GetServiceBEndpoint() + "/weather")
	if err != nil {
		err = errors.Wrap(err, "failed to parse integration url")
		telemetry.NotifyError(ctx, err)
		return nil, err
	}

	q := e.Query()
	q.Set("city", city)
	e.RawQuery = q.Encode()

	resp, err := httpclient.GetRequest(ctx, e.String())
	if err != nil {
		err = errors.Wrap(err, "failed to call service b")
		telemetry.NotifyError(ctx, err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var w Weather

	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		err = errors.Wrap(err, "failed to decode service b response")
		return nil, err
	}

	return &w, nil
}
