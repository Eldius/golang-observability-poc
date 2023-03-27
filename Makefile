
all-down: jaeger-down opensearch-down service-a-down
	@echo "everything is down"

network-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml down

network: network-down
	cd docker-environment ; docker compose -f docker-compose-network.yml up

jaeger-down:
	cd docker-environment ; docker compose -f docker-compose-jaeger.yml down

jaeger: jaeger-down network
	cd docker-environment ; docker compose -f docker-compose-jaeger.yml up

# opensearch-down:
# 	cd docker-environment ; docker compose -f docker-compose-opensearch.yml down

opensearch-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml down

# opensearch: network opensearch-down
# 	cd docker-environment ; docker compose -f docker-compose-opensearch.yml up

opensearch: opensearch-down
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml up

service-a-down:
	cd docker-environment ; docker compose -f docker-compose-network.yml -f docker-compose-opensearch.yml -f docker-compose-service-a.yml down

service-a: service-a-down
	$(eval PREPPER_ENDPOINT := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' data-prepper))
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	cd docker-environment ; PREPPER_ENDPOINT=$(PREPPER_ENDPOINT) docker compose \
		-f docker-compose-service-a.yml \
		-f docker-compose-network.yml \
		-f docker-compose-opensearch.yml \
		up \
			--build

service-a-local:
	$(eval PREPPER_ENDPOINT := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' data-prepper))
	$(eval DB_HOST := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres))
	@echo "data prepper: $(PREPPER_ENDPOINT)"
	@echo "db:     $(DB_HOST)"
	cd rest-service-a ;
	  API_LOG_LEVEL="trace" \
		API_OTEL_ENDPOINT="$(PREPPER_ENDPOINT):21890" \
		API_DB_HOST=$(DB_HOST) \
		API_DB_PASS="P@ss" \
		API_TELEMETRY_REST_ENABLE=true \
		API_TELEMETRY_DB_ENABLE=true \
		go run ./cmd \
			--config rest-api-config.yaml
