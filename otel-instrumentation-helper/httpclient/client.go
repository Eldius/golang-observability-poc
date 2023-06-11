package httpclient

import (
	"context"
	"fmt"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

var defaultClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

type options struct {
	headers map[string]interface{}
}

type Option func(o *options) *options

func WithHeader(k, v string) Option {
	return func(o *options) *options {
		o.headers[k] = v
		return o
	}
}

func MakeRequest(ctx context.Context, url, method string, opts ...Option) (*http.Response, error) {
	opt := &options{headers: make(map[string]interface{})}
	for _, o := range opts {
		opt = o(opt)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	for k, v := range opt.headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	l := logger.GetLogger(ctx).WithFields(logrus.Fields{
		"url":     url,
		"method":  req.Method,
		"headers": opt.headers,
	})

	l.Debug("preparing to make a request")

	res, err := defaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	l.Debug("request made")

	return res, nil
}

func GetRequest(ctx context.Context, url string, opts ...Option) (*http.Response, error) {
	return MakeRequest(ctx, url, http.MethodGet, opts...)
}

func GetHTTPClient() *http.Client {
	return defaultClient
}
