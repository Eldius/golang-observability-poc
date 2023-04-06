
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
	go get -u github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper@$(shell git rev-parse --short HEAD)


docker-down:
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
			eldius/service-a:dev

docker-up-opensearch: docker-down build-docker

	docker run \
		--rm \
		--name service_a \
		--network opensearch_default \
		-d \
		-m 16m \
		-p 8080:8080 \
		-e "API_OTEL_TRACE_ENDPOINT=data-prepper:21890" \
		-e "API_OTEL_METRICS_ENDPOINT=data-prepper:21891" \
		-e "API_DB_HOST=postgres" \
		-e "API_DB_PASS=P@ss" \
		-e "API_TELEMETRY_REST_ENABLE=true" \
		-e "API_TELEMETRY_DB_ENABLE=true" \
		-e "API_LOG_LEVEL=trace" \
			eldius/service-a:dev

lint:
	golangci-lint run --config=~/.golangci.yaml ./...
