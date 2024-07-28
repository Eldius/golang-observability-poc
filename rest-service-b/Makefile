
OBSERVABILITY_HOST ?= 192.168.0.68

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
	@echo "PWD: $(PWD)"
	docker buildx build \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		--no-cache \
		-t eldius/service-b:dev \
		-t eldius/service-b:$(COMMIT_HASH) \
			.

tidy:
	go mod tidy

build: tidy
	go build \
		-o .bin/app \
		-v \
		-a \
		-ldflags '-extldflags "-static"' \
		-ldflags "-X \"github.com/eldius/rest-api/rest-service-b/internal/config.Version=$(shell git rev-parse --short HEAD )\"" \
		./cmd/

update-library:
	go get -u github.com/eldius/golang-observability-poc/otel-instrumentation-helper@$(shell git rev-parse --short HEAD)

docker-up-opensearch: build-docker docker-down
	docker run \
		--rm \
		--name service_b \
		--network opensearch_default \
		-d \
		-m 16m \
		--log-driver=fluentd \
		--log-opt fluentd-address=localhost:24224 \
		-e "API_OTEL_TRACE_ENDPOINT=data-prepper:21890" \
		-e "API_OTEL_METRICS_ENDPOINT=data-prepper:21891" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-b:dev

docker-up-opensearch-k8s: docker-down build-docker
	$(eval FLUENTBIT_HOST := $(shell ./../fetch_ports.sh fluent-bit 24224 observability))
	$(eval COLLECTOR_TRACE_HOST := $(shell ./../fetch_ports.sh otel-collector 55689 observability))
	$(eval COLLECTOR_METRICS_HOST := $(shell ./../fetch_ports.sh otel-collector 55690 observability))
	docker run \
		--rm \
		--name service_b \
		-d \
		-m 16m \
		--network services_network \
		--log-driver=fluentd \
		--log-opt fluentd-address=$(FLUENTBIT_HOST) \
		-e "API_OTEL_TRACE_ENDPOINT=$(COLLECTOR_TRACE_HOST)" \
		-e "API_OTEL_METRICS_ENDPOINT=$(COLLECTOR_METRICS_HOST)" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
		-e "API_INTEGRATION_WEATHER_APIKEY=$(WEATHER_APIKEY)" \
			eldius/service-b:dev

docker-up-grafana: docker-down build-docker
	@echo "WEATHER_APIKEY: $(WEATHER_APIKEY)"
	docker run \
		--rm \
		--name service_b \
		--network services_network \
		-d \
		-m 16m \
		--log-driver=loki \
		--log-opt loki-url="http://$(OBSERVABILITY_HOST):3100/loki/api/v1/push" \
		--log-opt loki-external-labels="container_name={{.Name}},cluster=compose,container_id={{.ID}}" \
		-e "API_OTEL_TRACE_ENDPOINT=$(OBSERVABILITY_HOST):55689" \
		-e "API_OTEL_METRICS_ENDPOINT=$(OBSERVABILITY_HOST):55690" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
		-e "API_INTEGRATION_WEATHER_APIKEY=$(WEATHER_APIKEY)" \
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
		-e "API_INTEGRATION_WEATHER_APIKEY=$(WEATHER_APIKEY)" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-b:dev


docker-up-jaeger-k8s: docker-down build-docker
	docker run \
		--rm \
		--name service_b \
		-d \
		-m 16m \
		--network services_network \
		--log-driver=fluentd \
		--log-opt fluentd-address=$(CLUSTER_IP):24224 \
		-e "API_OTEL_TRACE_ENDPOINT=$(CLUSTER_IP):4317" \
		-e "API_OTEL_METRICS_ENDPOINT=$(CLUSTER_IP):21891" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
		-e "API_INTEGRATION_WEATHER_APIKEY=$(WEATHER_APIKEY)" \
			eldius/service-b:dev

lint:
	golangci-lint  run --config=../.golangci.yaml ./...
