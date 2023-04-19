
OPENSEARCH_IP := 192.168.0.36

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

services-network:
	-docker network create services_network

services-network-down:
	-docker network rm services_network

services-opensearch-k8s: $(APIS) services-network
	@echo "Services starting - Opensearch..."
	for dir in $(APIS); do \
		$(MAKE) -C $$dir docker-up-opensearch-k8s; \
	done

services-jaeger-k8s: $(APIS) services-network
	@echo "Services starting - Jaeger..."
	for dir in $(APIS); do \
		$(MAKE) -C $$dir docker-up-jaeger-k8s; \
	done

services-jaeger:
	@echo "Services starting - Jaeger..."
	@echo "Services starting - Opensearch..."
	for dir in $(APIS); do \
		$(MAKE) -C $$dir docker-up-jaeger; \
	done

services-down:
	@echo "Services stopping..."
	for dir in $(APIS); do \
		$(MAKE) -C $$dir docker-down; \
	done

tidy: $(APPS)
	for dir in $(APPS); do \
		$(MAKE) -C $$dir tidy; \
	done

lint: $(APPS)
	for dir in $(APPS); do \
		echo "linting $$dir..."; \
		$(MAKE) -C $$dir lint; \
	done

update-library:
	$(eval CURR_DIR := $(PWD))
	$(MAKE) -C apps/rest-service-a update-library
	$(MAKE) -C apps/rest-service-b update-library

weather:
	http http://localhost:8080/weather city=="Rio de Janeiro"

watch-service-a:
	watch -n 10 'curl -i localhost:8080/weather?city=Rio%20de%20Janeiro -H "Authorization: 854bf4f2-cb7d-11ed-bf82-00155d485640"'

exporting:
	docker run \
		--name data-prepper \
		--rm \
		-p 4900:4900 \
		-v ${PWD}/docker-environment/opensearch/configs/data:/usr/share/data-prepper/data \
		-v ${PWD}/docker-environment/opensearch/configs/logstash.conf:/usr/share/data-prepper/pipelines/pipelines.conf opensearchproject/data-prepper:latest

test-logs:
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
			eldius/service-a:dev

test:
	kubectl \
		run test-alpine \
		--image=alpine \
		--env="PS1='[\u@\h \W]\$ '" \
		-i \
		--tty \
		--restart=Never \
		--command -- sh

ks-terraform-opensearch-apply: ks-terraform-opensearch-init
	cd k3s-environment/terraform ; ELASTICSEARCH_USERNAME=admin \
		ELASTICSEARCH_PASSWORD=admin \
		ELASTICSEARCH_URL=https://$(OPENSEARCH_IP):9200 \
			TF_LOG=debug terraform apply

ks-terraform-opensearch-init:
	cd k3s-environment/terraform ; ELASTICSEARCH_USERNAME=admin \
		ELASTICSEARCH_PASSWORD=admin \
		ELASTICSEARCH_URL=https://$(OPENSEARCH_IP):9200 \
			terraform init

ks-terraform-opensearch-destroy: ks-terraform-opensearch-init
	cd k3s-environment/terraform ; ELASTICSEARCH_USERNAME=admin \
		ELASTICSEARCH_PASSWORD=admin \
		ELASTICSEARCH_URL=https://$(OPENSEARCH_IP):9200 \
			terraform destroy

jaeger-test:
	docker run \
		--rm \
		--name jaeger-quyery \
		-m 16m \
		-p 16687:16687 \
		-p 16686:16686 \
		-e SPAN_STORAGE_TYPE=elasticsearch \
		-e ES_SERVER_URLS=https://$(OPENSEARCH_IP):9200 \
		--log-driver=fluentd \
		--log-opt fluentd-address=192.168.0.36:24224 \
		-v "$(PWD)/docker-environment/opensearch/configs/root-ca.pem:/root-ca.pem:ro" \
			jaegertracing/jaeger-query:latest \
				--es.tls.skip-host-verify \
				--es.tls.ca "/root-ca.pem" \
				--admin.http.host-port ":16687" \
				--es.tls.enabled \
				--es.username admin \
				--es.password admin

	# docker run -d --rm \
	# -p 16686:16686 \
	# -p 16687:16687 \
	# -e SPAN_STORAGE_TYPE=elasticsearch \
	# -e ES_SERVER_URLS=https://$(OPENSEARCH_IP):9200 \
	# jaegertracing/jaeger-query:1.18

#jaeger-query
