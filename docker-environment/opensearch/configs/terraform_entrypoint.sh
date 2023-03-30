#!/bin/sh

# apk add curl
# curl -i -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}"  "http://kibana:5601/api/status"
# curl -i -u admin:admin  "http://${ELASTICSEARCH_URL}:9200/_cluster/health"

# wget --spider "http://${ELASTICSEARCH_URL}:9200/_cluster/health"

rm *.tfstate*

terraform init && terraform apply -auto-approve
