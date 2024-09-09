# Kind Kubernetes Cluster

This project is designed to present and test my Proof of Concepts (PoCs) and studies related to Kubernetes tools and configurations.

## Prerequisites

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [Helm](https://helm.sh/)

## Overview

This project leverages **Kind** (Kubernetes IN Docker) to create a local Kubernetes cluster for testing various tools like ArgoCD, Prometheus, Zabbix, and more. All operations can be executed via the Makefile.

## Commands Overview

### 1. Creating a Custom Cluster

To create a local Kubernetes cluster using Kind with a predefined configuration:

```bash
make create
```

This command uses the `config.yaml` file to provision the cluster with specific configurations, such as node settings, port mappings, etc.

### 2. Deleting the Cluster

To remove the existing Kubernetes cluster:

```bash
make delete
```

This command deletes the Kind cluster, cleaning up all the resources created in the process.

### 3. Creating ArgoCD

To install **ArgoCD** into the Kubernetes cluster:

```bash
make create-argocd
```

This installs ArgoCD using Helm with the custom values provided in `charts/argocd/values.yaml`, creating the `argocd` namespace if it doesn't exist.

### 4. Upgrading ArgoCD

To upgrade the existing ArgoCD installation:

```bash
make upgrade-argocd
```

This command upgrades the ArgoCD installation with the latest version or any changes specified in the `charts/argocd/values.yaml`.

### 5. Deleting ArgoCD

To uninstall ArgoCD from the cluster:

```bash
make delete-argocd
```

This removes ArgoCD and deletes the `argocd` namespace from the cluster.

### 6. Creating the Metrics Server

To install the **Metrics Server**:

```bash
make create-metrics-server
```

This deploys the Metrics Server using Helm and custom values from `charts/metrics-server/values.yaml` into the `kube-system` namespace.

### 7. Upgrading the Metrics Server

To upgrade the Metrics Server:

```bash
make upgrade-metrics-server
```

This command upgrades the Metrics Server with any new configurations or updates specified in the `charts/metrics-server/values.yaml`.

### 8. Deleting the Metrics Server

To uninstall the Metrics Server:

```bash
make delete-metrics-server
```

This command removes the Metrics Server from the `kube-system` namespace.

### 9. Creating Prometheus Stack

To install the **Prometheus Stack** for monitoring:

```bash
make create-prometheus-stack
```

This installs the Prometheus Stack (Prometheus, Alertmanager, Grafana, etc.) using Helm with the values from `charts/prometheus-stack/values.yaml`, creating the `prometheus` namespace.

### 10. Upgrading Prometheus Stack

To upgrade the Prometheus Stack:

```bash
make upgrade-prometheus-stack
```

This upgrades the Prometheus Stack with any changes or updates specified in `charts/prometheus-stack/values.yaml`.

### 11. Deleting Prometheus Stack

To uninstall the Prometheus Stack:

```bash
make delete-prometheus-stack
```

This command removes Prometheus and deletes the `prometheus` namespace.

### 12. Creating Jaeger

To deploy **Jaeger** for distributed tracing:

```bash
make create-jaeger
```

This applies the Jaeger deployment from the `manifests/jaeger/values.yaml`.

### 13. Deleting Jaeger

To delete the Jaeger deployment:

```bash
make delete-jaeger
```

This removes the Jaeger deployment.

### 14. Creating Odigos

To install **Odigos**:

```bash
make create-odigos
```

This command applies the Odigos deployment defined in `labs/odigos/values.yaml`.

### 15. Deleting Odigos

To delete the Odigos deployment:

```bash
make delete-odigos
```

This command removes the Odigos deployment from the cluster.

### 16. Creating Prometheus Adapter

To install the **Prometheus Adapter** for custom metrics:

```bash
make create-prometheus-adapter
```

This installs the Prometheus Adapter using Helm with the values from `charts/prometheus-adapter/values.yaml`, creating the `custom-metrics` namespace.

### 17. Deleting Prometheus Adapter

To uninstall the Prometheus Adapter:

```bash
make delete-prometheus-adapter
```

This command removes the Prometheus Adapter and deletes the `custom-metrics` namespace.

### 18. Creating Zabbix

To install **Zabbix**:

```bash
make create-zabbix
```

This installs Zabbix using Helm with the values from `charts/zabbix/values.yaml`, creating the `zabbix` namespace.

### 19. Deleting Zabbix

To uninstall Zabbix:

```bash
make delete-zabbix
```

This removes Zabbix and deletes the `zabbix` namespace.

### 20. Creating Harbor

To install **Harbor** for container registry:

```bash
make create-harbor
```

This installs Harbor using Helm with the values from `charts/harbor/values.yaml`, creating the `harbor` namespace.

### 21. Deleting Harbor

To uninstall Harbor:

```bash
make delete-harbor
```

This removes Harbor and deletes the `harbor` namespace.
