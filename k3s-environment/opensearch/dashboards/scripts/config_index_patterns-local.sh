#!/bin/sh
while ! curl -i \
    --fail \
    --insecure \
    -X POST \
    -u 'admin:admin' \
    'http://opensearch-dashboards:5601/api/saved_objects/index-pattern/application-logs' \
    -H 'osd-xsrf: true' \
    -H 'Content-Type: application/json' \
    -H 'securitytenant: global' \
    -d '{
        "attributes": {
            "title": "application-logs",
            "timeFieldName": "@timestamp"
        }
    }'
do
    sleep 5
done

curl -i \
        --fail \
        --insecure \
        -X POST \
        -u 'admin:admin' \
        'http://opensearch-dashboards:5601/api/saved_objects/index-pattern/metrics-otel-v1-*' \
        -H 'osd-xsrf: true' \
        -H 'Content-Type: application/json' \
        -H 'securitytenant: global' \
        -d '{
            "attributes": {
                "title": "metrics-otel-v1-*",
                "timeFieldName": "@timestamp"
            }
        }'
