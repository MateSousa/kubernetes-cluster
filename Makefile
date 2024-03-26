.PHONY: create delete cilium jaeger prometheus istio

create:
	kind create cluster --config config.yaml

delete: 
	kind delete cluster

create-argocd:
	helm install argocd argo/argo-cd -n argocd --create-namespace

upgrade-argocd:
	helm upgrade argocd argo/argo-cd -n argocd

delete-argocd:
	helm uninstall argocd -n argocd
