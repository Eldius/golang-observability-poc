package weather

import (
	"context"
	"encoding/json"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/config"
	"net/http"
	"net/url"
)

func GetWeather(ctx context.Context, city string) (*Weather, error) {
	ctx, closer := telemetry.StartSpan(ctx, "GetWeatherIntegration")
	defer closer()

	telemetry.AddAttribute(ctx, "city", city)

	endpoint, err := url.Parse(config.GetWeatherServiceEndpoint())
	if err != nil {
		return nil, err
	}
	qp := endpoint.Query()
	qp.Set("city", city)

	endpoint.RawQuery = qp.Encode()

	c := http.Client{}

	resp, err := c.Get(endpoint.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var w Weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}

	return &w, nil
}
