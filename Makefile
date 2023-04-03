
APPS := $(wildcard apps/*/.)

APIS := $(wildcard apps/rest-service*/.)

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

services-opensearch: $(APIS)
	@echo "Services starting - Opensearch..."
	for dir in $(APIS); do \
		$(MAKE) -C $$dir docker-up-opensearch; \
	done
	# @(MAKE) -C $(PWD)/apps/rest-service-a docker-up-opensearch
	# @(MAKE) -C $(PWD)/rest-service-b docker-up-opensearch

services-jaeger:
	@echo "Services starting - Jaeger..."
	@(MAKE) -C apps/rest-service-a docker-up-jaeger
	@(MAKE) -C apps/rest-service-b docker-up-jaeger

services-down: service-a-down service-b-down
	@echo "Services stopping..."
	@(MAKE) -C apps/rest-service-a docker-down
	@(MAKE) -C apps/rest-service-b docker-down

tidy: $(APPS)
	for dir in $(APPS); do \
		$(MAKE) -C $$dir tidy; \
	done

update-library:
	$(eval CURR_DIR := $(PWD))
	$(MAKE) -C apps/rest-service-a update-library
	$(MAKE) -C apps/rest-service-b update-library

weather:
	http http://localhost:8080/weather city=="Rio de Janeiro"

watch-service-a:
	watch -n 10 'curl -i localhost:8080/weather?city=Rio%20de%20Janeiro -H "Authorization: 854bf4f2-cb7d-11ed-bf82-00155d485640"'
