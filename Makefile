.PHONY: create delete create-argocd upgrade-argocd delete-argocd create-metrics-server upgrade-metrics-server delete-metrics-server create-prometheus-stack upgrade-prometheus-stack delete-prometheus-stack create-jaeger delete-jaeger create-odigos delete-odigos create-prometheus-adapter delete-prometheus-adapter create-zabbix delete-zabbix create-harbor delete-harbor

create:
	kind create cluster --config config.yaml

delete: 
	kind delete cluster

create-argocd:
	helm install argocd argo/argo-cd -n argocd --create-namespace -f charts/argocd/values.yaml

upgrade-argocd:
	helm upgrade argocd argo/argo-cd -n argocd -f charts/argocd/values.yaml

delete-argocd:
	helm uninstall argocd -n argocd && kubectl delete namespace argocd

create-metrics-server:
	helm install metrics-server metrics-server/metrics-server -n kube-system -f charts/metrics-server/values.yaml

upgrade-metrics-server:
	helm upgrade metrics-server metrics-server -n kube-system -f charts/metrics-server/values.yaml

delete-metrics-server:
	helm uninstall metrics-server -n kube-system

create-prometheus-stack:
	helm install prometheus-stack prometheus-community/kube-prometheus-stack -n prometheus --create-namespace -f charts/prometheus-stack/values.yaml

upgrade-prometheus-stack:
	helm upgrade prometheus-stack prometheus-community/kube-prometheus-stack -n prometheus -f charts/prometheus-stack/values.yaml

delete-prometheus-stack:
	helm uninstall prometheus-stack -n prometheus && kubectl delete namespace prometheus

create-jaeger:
	kubectl apply -f manifests/jaeger/values.yaml

delete-jaeger:
	kubectl delete -f manifests/jaeger/values.yaml

create-odigos:
	kubectl apply -f labs/odigos/values.yaml

delete-odigos:
	kubectl delete -f labs/odigos/values.yaml

create-prometheus-adapter:
	helm install prometheus-adapter prometheus-community/prometheus-adapter -n custom-metrics --create-namespace -f charts/prometheus-adapter/values.yaml

delete-prometheus-adapter:
	helm uninstall prometheus-adapter -n custom-metrics && kubectl delete namespace custom-metrics

create-zabbix:
	helm install my-zabbix-test zabbix-community/zabbix --version 5.0.1 -n zabbix --create-namespace -f charts/zabbix/values.yaml

delete-zabbix:
	helm uninstall my-zabbix-test -n zabbix && kubectl delete namespace zabbix

create-harbor:
	helm install poc-harbor harbor/harbor -n harbor --create-namespace -f charts/harbor/values.yaml

delete-harbor:
	helm uninstall poc-harbor -n harbor && kubectl delete namespace harbor
