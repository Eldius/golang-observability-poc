package httpclient

import (
    "context"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "net/http"
)

var (
    defaultClient = &http.Client{
        Transport: otelhttp.NewTransport(http.DefaultTransport),
    }
)

type options struct {
    headers map[string]interface{}
}

type Option func(o *options) *options

func MakeRequest(ctx context.Context, url string, method string, opts ...Option) (*http.Response, error) {
    opt := &options{}
    for _, o := range opts {
        opt = o(opt)
    }
    req, err := http.NewRequestWithContext(ctx, method, url, nil)
    if err != nil {
        return nil, err
    }

    res, err := defaultClient.Do(req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func GetRequest(ctx context.Context, url string, opts ...Option) (*http.Response, error) {
    return MakeRequest(ctx, url, http.MethodGet, opts...)
}

func GetHttpClient() *http.Client {
    return defaultClient
}
