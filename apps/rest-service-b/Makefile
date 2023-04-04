
start:
	$(eval OTEL_ENDPOINT := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' jaeger))
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	@echo "jaeger: $(OTEL_ENDPOINT)"
	@echo "db:     $(DB_HOST)"
	API_LOG_LEVEL="trace" \
		API_OTEL_ENDPOINT="$(OTEL_ENDPOINT):4317" \
		API_DB_HOST=$(DB_HOST) \
		API_DB_PASS="P@ss" \
		go run ./cmd \
			--config ./rest-api-config.yaml

build-docker:
	$(eval COMMIT_HASH := $(shell git rev-parse --short HEAD))
	@echo "COMMIT_HASH: $(COMMIT_HASH)"
	docker build \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		-t eldius/service-b:dev \
			.

tidy:
	go mod tidy

build: tidy
	go build \
		-o .bin/app \
		-v \
		-a \
		-ldflags '-extldflags "-static"' \
		-ldflags "-X \"github.com/eldius/rest-api/rest-service-a/internal/config.Version=$(shell git rev-parse --short HEAD )\"" \
		./cmd/

update-library:
	go get -u github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper@$(shell git rev-parse --short HEAD)

docker-up-opensearch: build-docker docker-down
	docker run \
		--rm \
		--name service_b \
		--network opensearch_default \
		-d \
		-m 16m \
		-e "API_OTEL_TRACE_ENDPOINT=data-prepper:21890" \
		-e "API_OTEL_METRICS_ENDPOINT=data-prepper:21891" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
			eldius/service-b:dev

docker-down:
	-docker kill service_b
	-docker rm service_b

docker-up-jaeger: docker-down build-docker
	docker run \
		--rm \
		--name service_b \
		--network jaeger_default \
		-d \
		-m 16m \
		-e "API_OTEL_TRACE_ENDPOINT=jaeger:4317" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_INTEGRATION_WEATHER_APIKEY=$(WEATHER_API_KEY)" \
			eldius/service-b:dev
