.PHONY: create delete create-argocd upgrade-argocd delete-argocd create-metrics-server upgrade-metrics-server delte-metrics-server

create:
	kind create cluster --config config.yaml

delete: 
	kind delete cluster

create-argocd:
	helm install argocd argo/argo-cd -n argocd --create-namespace -f charts/argocd/values.yaml

upgrade-argocd:
	helm upgrade argocd argo/argo-cd -n argocd -f charts/argocd/values.yaml

delete-argocd:
	helm uninstall argocd -n argocd --delete-namespace

create-metrics-server:
	helm install metrics-server metrics-server -n kube-system -f charts/metrics-server/values.yaml

upgrade-metrics-server:
	helm upgrade metrics-server metrics-server -n kube-system -f charts/metrics-server/values.yaml

delete-metrics-server:
	helm uninstall metrics-server -n kube-system

create-prometheus-stack:
	helm install prometheus-stack prometheus-community/kube-prometheus-stack -n prometheus --create-namespace -f charts/prometheus-stack/values.yaml

upgrade-prometheus-stack:
	helm upgrade prometheus-stack prometheus-community/kube-prometheus-stack -n prometheus -f charts/prometheus-stack/values.yaml

delete-prometheus-stack:
	helm uninstall prometheus-stack -n prometheus --delete-namespace

create-jaeger:
	kubectl apply -f manifests/jaeger/values.yaml

delete-jaeger:
	kubectl delete -f manifests/jaeger/values.yaml

create-odigos:
	kubectl apply -f labs/odigos/values.yaml

delete-odigos:
	kubectl delete -f labs/odigos/values.yaml
