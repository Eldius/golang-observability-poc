
all-down: jaeger-down opensearch-down service-a-down
	@echo "everything is down"

network-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml down

network: network-down
	cd docker-environment ; docker compose -f docker-compose-network.yml up

jaeger-down:
	cd docker-environment ; docker compose -f docker-compose-jaeger.yml down

jaeger: jaeger-down
	cd docker-environment ; docker compose -f docker-compose-jaeger.yml up

# opensearch-down:
# 	cd docker-environment ; docker compose -f docker-compose-opensearch.yml down

opensearch-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml down

# opensearch: network opensearch-down
# 	cd docker-environment ; docker compose -f docker-compose-opensearch.yml up

opensearch: opensearch-down
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml up

service-a-opensearch-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml down

service-a-opensearch: service-a-opensearch-down
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	cd docker-environment ; OTEL_ENDPOINT=data-prepper docker compose \
		-f docker-compose-service-a.yml \
		-f docker-compose-network.yml \
		-f docker-compose-opensearch.yml \
		up \
			--build

service-a-jaeger-down:
	cd docker-environment ; docker compose \
		-f docker-compose-service-a-jaeger.yml \
			down

service-a-jaeger: service-a-jaeger-down
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	cd docker-environment ; docker compose \
		-f docker-compose-service-a-jaeger.yml \
		up \
			--build

service-a-db-down:
	cd docker-environment ; docker compose \
		-f docker-compose-db.yml \
		down

service-a-db: service-a-db-down
	cd docker-environment ; docker compose \
		-f docker-compose-db.yml \
		up


service-a-local:
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

testing:
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	cd docker-environment ; OTEL_ENDPOINT=jaeger:4317 docker compose \
		-f docker-compose-jaeger.yml \
		-f docker-compose-db.yml \
		up \
			--build
