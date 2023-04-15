
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

ks-fluent-bit: ks-fluent-bit-configmap
	cd opensearch/fluent-bit; kubectl apply -f .

ks-fluent-bit-down: ks-fluent-bit-configmap-down
	cd opensearch/fluent-bit; kubectl delete -f .

ks-opensearch-down:
	cd opensearch/opensearch; kubectl delete -f .

ks-data-prepper-configmap:
	kubectl create configmap data-prepper-config-files --from-file=opensearch/data-prepper/configs

ks-data-prepper-configmap-down:
	kubectl delete configmap data-prepper-config-files

ks-fluent-bit-configmap:
	kubectl create configmap fluent-bit-config-files --from-file=opensearch/fluent-bit/configs

ks-fluent-bit-configmap-down:
	kubectl delete configmap fluent-bit-config-files
