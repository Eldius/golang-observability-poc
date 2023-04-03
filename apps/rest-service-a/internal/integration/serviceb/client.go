package serviceb

import (
    "context"
    "encoding/json"
    "github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/httpclient"
    "github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
    "github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
    "net/url"
)

func GetWeather(ctx context.Context, city string) (*Weather, error) {

    e, err := url.Parse(config.GetServiceBEndpoint() + "/weather")
    if err != nil {
        telemetry.NotifyError(ctx, err)
        return nil, err
    }

    q := e.Query()
    q.Set("city", city)
    e.RawQuery = q.Encode()

    resp, err := httpclient.GetRequest(ctx, e.String())
    if err != nil {
        telemetry.NotifyError(ctx, err)
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
