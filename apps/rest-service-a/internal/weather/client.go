package weather

import (
	"encoding/json"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
)

func GetWeather(city string) (*Weather, error) {

	c := telemetry.GetHttpClient()

	resp, err := c.Get(config.GetServiceBEndpoint())
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
