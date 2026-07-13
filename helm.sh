#!/usr/bin/env bash

# chmod +x helm.sh
# ./helm.sh

mkdir -p helm/inference-orchestrator/templates/{api,redis,worker-classification,worker-regression,monitoring}

touch helm/inference-orchestrator/Chart.yaml
touch helm/inference-orchestrator/values.yaml
touch helm/inference-orchestrator/values-cloud.yaml

touch helm/inference-orchestrator/templates/api/deployment.yaml
touch helm/inference-orchestrator/templates/api/service.yaml

touch helm/inference-orchestrator/templates/redis/deployment.yaml
touch helm/inference-orchestrator/templates/redis/service.yaml
touch helm/inference-orchestrator/templates/redis/pvc.yaml

touch helm/inference-orchestrator/templates/worker-classification/deployment.yaml
touch helm/inference-orchestrator/templates/worker-classification/configmap.yaml
touch helm/inference-orchestrator/templates/worker-classification/scaledobject.yaml

touch helm/inference-orchestrator/templates/worker-regression/deployment.yaml
touch helm/inference-orchestrator/templates/worker-regression/configmap.yaml
touch helm/inference-orchestrator/templates/worker-regression/scaledobject.yaml

touch helm/inference-orchestrator/templates/monitoring/prometheus-deployment.yaml
touch helm/inference-orchestrator/templates/monitoring/prometheus-service.yaml
touch helm/inference-orchestrator/templates/monitoring/grafana-deployment.yaml
touch helm/inference-orchestrator/templates/monitoring/grafana-service.yaml
touch helm/inference-orchestrator/templates/monitoring/loki-deployment.yaml
touch helm/inference-orchestrator/templates/monitoring/loki-service.yaml
touch helm/inference-orchestrator/templates/monitoring/loki-pvc.yaml
touch helm/inference-orchestrator/templates/monitoring/promtail-daemonset.yaml
touch helm/inference-orchestrator/templates/monitoring/pushgateway-deployment.yaml
touch helm/inference-orchestrator/templates/monitoring/pushgateway-service.yaml