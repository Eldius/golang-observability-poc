#!/bin/sh

# apk add curl
# curl -i -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}"  "http://kibana:5601/api/status"
# curl -i -u elastic:changeme  "http://elasticsearch:9200/_cluster/health"

rm *.tfstate*

ls -lha

terraform init && terraform apply -auto-approve
