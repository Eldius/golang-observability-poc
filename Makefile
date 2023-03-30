
env-jaeger-down:
	cd docker-environment/jaeger ; docker compose \
		-f docker-compose-jaeger.yml \
		-f ../docker-compose-db.yml \
		down

env-jaeger: env-jaeger-down
	cd docker-environment/jaeger ; docker compose \
		-f docker-compose-jaeger.yml \
		-f ../docker-compose-db.yml \
		up \
		-d \
			--build

env-opensearch-down:
	cd docker-environment/opensearch ; docker compose \
		-f docker-compose-opensearch.yml \
		-f ../docker-compose-db.yml \
		down

env-opensearch: env-opensearch-down
	cd docker-environment/opensearch ; docker compose \
		-f docker-compose-opensearch.yml \
		-f ../docker-compose-db.yml \
		up \
		-d \
			--build

filebeat: filebeat-down
	cd docker-environment/opensearch ; docker compose \
		-f docker-compose-filebeat.yml \
		-f ../docker-compose-db.yml \
		up \
		-d \
			--build

filebeat-down:
	cd docker-environment/opensearch ; docker compose \
		-f docker-compose-filebeat.yml \
		-f ../docker-compose-db.yml \
		down

service-a-build-docker:
	cd apps/rest-service-a && \
		$(MAKE) build-docker


service-a-build:
	cd apps/rest-service-a && \
		$(MAKE) build

service-a-jaeger: service-a-build-docker
	docker run \
		--rm \
		--name service_a \
		--network jaeger_default \
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

service-a:
	docker run \
		--rm \
		--name service_a \
		--network jaeger_default \
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

service-a-opensearch:
	cd apps/rest-service-a && \
		$(MAKE) build-docker

	docker run \
		--rm \
		--name service_a \
		--network opensearch_default \
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

service-a-local-jaeger:
	$(eval JAEGER_ENDPOINT := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' jaeger))
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	cd rest-service-a ; API_LOG_LEVEL="trace" \
		API_OTEL_ENDPOINT="$(JAEGER_ENDPOINT):4317" \
		API_DB_HOST=$(DB_HOST) \
		API_DB_PASS="P@ss" \
		API_TELEMETRY_REST_ENABLE=true \
		API_TELEMETRY_DB_ENABLE=true \
		go run ./cmd \
			--config rest-api-config.yaml

watch-service-a:
	watch -n 10 'curl -i localhost:8080/ping -H "Authorization: 854bf4f2-cb7d-11ed-bf82-00155d485640"'
