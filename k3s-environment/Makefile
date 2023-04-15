
ks-opensearch:
	cd opensearch/opensearch; kubectl apply -f .

ks-dashboards:
	cd opensearch/dashboards; kubectl apply -f .

ks-dashboards-init:
	./opensearch/dashboards/scripts/config_index_patterns.sh

ks-dashboards-down:
	cd opensearch/dashboards; kubectl delete -f .

ks-data-prepper: ks-data-prepper-configmap
	cd opensearch/data-prepper; kubectl apply -f .

ks-data-prepper-down: ks-data-prepper-configmap-down
	cd opensearch/data-prepper; kubectl delete -f .

ks-opensearch-down:
	cd opensearch/opensearch; kubectl delete -f .

ks-data-prepper-configmap:
	kubectl create configmap data-prepper-config-files --from-file=opensearch/data-prepper/configs

ks-data-prepper-configmap-down:
	kubectl delete configmap data-prepper-config-files
