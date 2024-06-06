
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
	docker buildx build \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		-t eldius/service-a:dev \
			.

tidy:
	go mod tidy

build: tidy
	go build \
		-o .bin/app \
		-v \
		-a \
		-ldflags '-extldflags "-static"' \
		-ldflags "-X 'github.com/eldius/rest-api/rest-service-a/internal/config.Version=$(shell git rev-parse --short HEAD )'" \
		./cmd/

update-library:
	go get -u github.com/eldius/golang-observability-poc/otel-instrumentation-helper@$(shell git rev-parse --short HEAD)


docker-down:
	-docker kill db
	-docker rm db
	-docker kill service_a
	-docker rm service_a

docker-up-jaeger: docker-down build-docker
	docker run \
		--rm \
		--name service_a \
		--network jaeger_default \
		-d \
		-m 16m \
		-p 8080:8080 \
		-e "API_OTEL_TRACE_ENDPOINT=jaeger:4317" \
		-e "API_OTEL_METRICS_ENDPOINT=jaeger:4317" \
		-e "API_DB_HOST=postgres" \
		-e "API_DB_PASS=P@ss" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-a:dev

docker-up-opensearch: docker-down build-docker

	docker run \
		--rm \
		--name service_a \
		--network opensearch_default \
		-d \
		-m 16m \
		-p 8080:8080 \
		--log-driver=fluentd \
		--log-opt fluentd-address=localhost:24224 \
		-e "API_OTEL_TRACE_ENDPOINT=data-prepper:21890" \
		-e "API_OTEL_METRICS_ENDPOINT=data-prepper:21891" \
		-e "API_DB_HOST=postgres" \
		-e "API_DB_PASS=P@ss" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-a:dev

docker-up-opensearch-k8s: docker-down build-docker
	$(eval FLUENTBIT_HOST := $(shell ./../fetch_ports.sh fluent-bit 24224 observability))
	$(eval COLLECTOR_TRACE_HOST := $(shell ./../fetch_ports.sh otel-collector 55689 observability))
	$(eval COLLECTOR_METRICS_HOST := $(shell ./../fetch_ports.sh otel-collector 55690 observability))
	$(eval DB_HOST := $(shell ./../fetch_ports.sh postgres 5432 databases))

	@echo "FLUENTBIT_HOST:        $(FLUENTBIT_HOST)"
	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
	@echo "DB_HOST:               $(DB_HOST)"

	docker run \
		--rm \
		--name service_a \
		--network services_network \
		-d \
		-m 16m \
		-p 8080:8080 \
		--log-driver=fluentd \
		--log-opt fluentd-address=$(FLUENTBIT_HOST) \
		-e "API_OTEL_TRACE_ENDPOINT=$(COLLECTOR_TRACE_HOST)" \
		-e "API_OTEL_METRICS_ENDPOINT=$(COLLECTOR_METRICS_HOST)" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_DB_HOST=$(DB_HOST)" \
		-e "API_DB_USER=rest_api" \
		-e "API_DB_PASS=P@ss" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-a:dev

docker-up-grafana: docker-down build-docker
	docker run \
		--rm \
		--name db \
		--network services_network \
		-d \
		-m 16m \
		-p 5432:5432 \
		-e "POSTGRES_USER=rest_api" \
		-e "POSTGRES_DB=rest_api" \
		-e "POSTGRES_PASSWORD=P@ss" \
		  	postgres

	docker run \
		--rm \
		--name service_a \
		--network services_network \
		-d \
		-m 16m \
		-p 8080:8080 \
		--log-driver=loki \
		--log-opt loki-url="http://192.168.0.68:3100/loki/api/v1/push" \
		--log-opt loki-external-labels="container_name={{.Name}},cluster=compose,container_id={{.ID}}" \
		-e "API_OTEL_TRACE_ENDPOINT=192.169.0.68:55689" \
		-e "API_OTEL_METRICS_ENDPOINT=192.169.0.68:55690" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_DB_HOST=db" \
		-e "API_DB_USER=rest_api" \
		-e "API_DB_PASS=P@ss" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-a:dev

docker-up-skywalking-k8s: docker-down build-docker
	$(eval FLUENTBIT_HOST := $(shell ./../fetch_ports.sh fluent-bit 24224 observability))
	$(eval COLLECTOR_TRACE_HOST := $(shell ./../fetch_ports.sh otel-collector 55689 observability))
	$(eval COLLECTOR_METRICS_HOST := $(shell ./../fetch_ports.sh otel-collector 55690 observability))
	docker run \
		--rm \
		--name service_a \
		--network services_network \
		-d \
		-m 16m \
		-p 8080:8080 \
		--log-driver=fluentd \
		--log-opt "fluentd-address=$(FLUENTBIT_HOST)" \
		-e "API_OTEL_TRACE_ENDPOINT=$(COLLECTOR_TRACE_HOST)" \
		-e "API_OTEL_METRICS_ENDPOINT=$(COLLECTOR_METRICS_HOST)" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_DB_HOST=192.168.100.183" \
		-e "API_DB_PORT=32277" \
		-e "API_DB_PASS=P@ss" \
		-e "API_LOG_FORMAT=json" \
			eldius/service-a:dev

lint:
	golangci-lint  run --config=../.golangci.yaml ./...


docker-up-jaeger-k8s: docker-down build-docker
	docker run \
		--rm \
		--name service_a \
		--network services_network \
		-d \
		-m 16m \
		-p 8080:8080 \
		--log-driver=fluentd \
		--log-opt fluentd-address=192.168.0.36:24224 \
		-e "API_OTEL_TRACE_ENDPOINT=192.168.0.36:4317" \
		-e "API_OTEL_METRICS_ENDPOINT=192.168.0.36:21891" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
		-e "API_DB_HOST=192.168.0.36" \
		-e "API_DB_PASS=P@ss" \
			eldius/service-a:dev
