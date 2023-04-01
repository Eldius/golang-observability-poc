
APPS := $(wildcard apps/*/.)

env-jaeger-down: services-down
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

env-opensearch-down: services-down
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

service-a-build-docker:
	cd apps/rest-service-a && \
		$(MAKE) build-docker


service-a-down:
	-docker kill service_a
	-docker rm service_a

service-a-jaeger: service-a-down service-a-build-docker
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

service-a-opensearch: service-a-down
	cd apps/rest-service-a && \
		$(MAKE) build-docker

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

service-b-build-docker:
	cd apps/rest-service-b && \
		$(MAKE) build-docker

service-b-opensearch: service-b-build-docker service-b-down
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

service-b-down:
	-docker kill service_b
	-docker rm service_b

service-b-jaeger: service-b-down service-b-build-docker
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
			eldius/service-b:dev



services-opensearch: service-b-opensearch service-a-opensearch
	@echo "Services started..."

services-jaeger: service-a-jaeger service-b-jaeger
	@echo "Services started..."

services-down: service-a-down service-b-down
	@echo "Services stopped..."

watch-service-a:
	watch -n 10 'curl -i localhost:8080/weather?city=Rio%20de%20Janeiro -H "Authorization: 854bf4f2-cb7d-11ed-bf82-00155d485640"'

tidy: $(APPS)
	for dir in $(APPS); do \
		$(MAKE) -C $$dir tidy; \
	done

update-library:
	$(eval CURR_DIR := $(PWD))
	$(MAKE) -C apps/rest-service-a update-library
	$(MAKE) -C apps/rest-service-b update-library
