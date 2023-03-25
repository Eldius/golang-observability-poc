
start:
	$(eval OTEL_ENDPOINT := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' jaeger))
	@echo "endpoint: $(OTEL_ENDPOINT)"
	API_LOG_LEVEL="trace" API_OTEL_ENDPOINT="$(OTEL_ENDPOINT):4317" go run ./cmd --config ./rest-api-config.yaml
