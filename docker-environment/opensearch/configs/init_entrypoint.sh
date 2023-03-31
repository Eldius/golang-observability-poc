#!/bin/sh

apk add curl

echo ""
echo "####################"
echo "## creating index ##"
echo "####################"
echo ""

curl -i \
    --insecure \
    -X PUT \
    -u 'admin:admin' \
    'https://node-0.example.com:9200/application-logs-00001' \
    -H 'Content-Type: application/json' \
    -d '{
        "settings": {
            "index": {
            "number_of_shards": 2,
            "number_of_replicas": 1
            }
        },
        "aliases": {
            "application-logs": {}
        }
    }' || exit 1


echo ""
echo ""
echo "############################"
echo "## creating index pattern ##"
echo "############################"
echo ""

curl -i \
    --insecure \
    -X POST \
    -u 'admin:admin' \
    "https://node-0.example.com:9200/_index_template/application-logs" \
    -H 'osd-xsrf: true' \
    -H 'Content-Type: application/json' \
    -d '{
        "index_patterns": [
            "application-logs"
        ]
    }' || exit 1
