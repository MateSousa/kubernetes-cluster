# Cluster management targets
create:
	kind create cluster --config config.yaml

delete:
	kind delete cluster

HELM_COMPONENTS_WITH_UPGRADE = argocd metrics-server prometheus-stack jaeger-stack kafka burrito
HELM_COMPONENTS_NO_UPGRADE = prometheus-adapter

HELM_COMPONENTS = $(HELM_COMPONENTS_WITH_UPGRADE) $(HELM_COMPONENTS_NO_UPGRADE)

argocd_RELEASE_NAME = argocd
argocd_CHART = argo/argo-cd
argocd_NAMESPACE = argocd
argocd_VALUES_FILE = charts/argocd/values.yaml
argocd_DELETE_NAMESPACE = yes

metrics-server_RELEASE_NAME = metrics-server
metrics-server_CHART = metrics-server/metrics-server
metrics-server_NAMESPACE = kube-system
metrics-server_VALUES_FILE = charts/metrics-server/values.yaml
metrics-server_DELETE_NAMESPACE = no

prometheus-stack_RELEASE_NAME = prometheus-stack
prometheus-stack_CHART = prometheus-community/kube-prometheus-stack
prometheus-stack_NAMESPACE = prometheus
prometheus-stack_VALUES_FILE = charts/prometheus-stack/values.yaml
prometheus-stack_DELETE_NAMESPACE = yes

prometheus-adapter_RELEASE_NAME = prometheus-adapter
prometheus-adapter_CHART = prometheus-community/prometheus-adapter
prometheus-adapter_NAMESPACE = custom-metrics
prometheus-adapter_VALUES_FILE = charts/prometheus-adapter/values.yaml
prometheus-adapter_DELETE_NAMESPACE = yes

jaeger-stack_RELEASE_NAME = jaeger
jaeger-stack_CHART = jaegertracing/jaeger
jaeger-stack_NAMESPACE = observability
jaeger-stack_VALUES_FILE = charts/jaeger/values.yaml
jaeger-stack_DELETE_NAMESPACE = no

kafka_RELEASE_NAME = kafka
kafka_CHART = oci://registry-1.docker.io/bitnamicharts/kafka
kafka_NAMESPACE = kafka
kafka_VALUES_FILE = charts/kafka/values.yaml
kafka_DELETE_NAMESPACE = no

burrito_RELEASE_NAME = burrito
burrito_CHART = oci://ghcr.io/padok-team/charts/burrito
burrito_NAMESPACE = burrito-system
burrito_VALUES_FILE = charts/burrito/values.yaml
burrito_DELETE_NAMESPACE = no

define HELM_TARGETS
create-$(1):
	helm install $$($(1)_RELEASE_NAME) $$($(1)_CHART) -n $$($(1)_NAMESPACE) --create-namespace -f $$($(1)_VALUES_FILE)

delete-$(1):
	helm uninstall $$($(1)_RELEASE_NAME) -n $$($(1)_NAMESPACE) $$(if $$(filter yes,$$($(1)_DELETE_NAMESPACE)),--delete-namespace)

ifeq ($(filter $(1),$(HELM_COMPONENTS_WITH_UPGRADE)), $(1))
upgrade-$(1):
	helm upgrade $$($(1)_RELEASE_NAME) $$($(1)_CHART) -n $$($(1)_NAMESPACE) -f $$($(1)_VALUES_FILE)
endif
endef

$(foreach component,$(HELM_COMPONENTS),$(eval $(call HELM_TARGETS,$(component))))

HELM_CREATE_TARGETS = $(foreach component,$(HELM_COMPONENTS),create-$(component))
HELM_DELETE_TARGETS = $(foreach component,$(HELM_COMPONENTS),delete-$(component))
HELM_UPGRADE_TARGETS = $(foreach component,$(HELM_COMPONENTS_WITH_UPGRADE),upgrade-$(component))

KUBECTL_COMPONENTS = jaeger odigos

jaeger_MANIFEST_FILE = manifests/jaeger/values.yaml
odigos_MANIFEST_FILE = labs/odigos/values.yaml

define KUBECTL_TARGETS
create-$(1):
	kubectl apply -f $$($(1)_MANIFEST_FILE)

delete-$(1):
	kubectl delete -f $$($(1)_MANIFEST_FILE)
endef

$(foreach component,$(KUBECTL_COMPONENTS),$(eval $(call KUBECTL_TARGETS,$(component))))

KUBECTL_CREATE_TARGETS = $(foreach component,$(KUBECTL_COMPONENTS),create-$(component))
KUBECTL_DELETE_TARGETS = $(foreach component,$(KUBECTL_COMPONENTS),delete-$(component))
KUBECTL_TARGETS = $(KUBECTL_CREATE_TARGETS) $(KUBECTL_DELETE_TARGETS)

.PHONY: create delete $(HELM_CREATE_TARGETS) $(HELM_UPGRADE_TARGETS) $(HELM_DELETE_TARGETS) $(KUBECTL_TARGETS)
