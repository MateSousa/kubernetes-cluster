.PHONY: create delete cilium jaeger prometheus istio

create:
	kind create cluster --config config.yaml

delete: 
	kind delete cluster
