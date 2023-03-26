
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
