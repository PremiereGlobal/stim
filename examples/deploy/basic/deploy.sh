#!/bin/sh

echo "Deploying Grafana to ${DEPLOY_ENVIRONMENT} in instance ${DEPLOY_INSTANCE} in cluster ${DEPLOY_CLUSTER}"

helm upgrade --install \
  --namespace ${NAMESPACE} \
  --version ${HELM_CHART_VERSION} \
  --set adminPassword=${ADMIN_PASS} \
  ${HELM_DEPLOYMENT_NAME} \
  stable/grafana
