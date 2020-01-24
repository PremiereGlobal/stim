#!/bin/bash
set -eu -o pipefail

# Exit if not running with `stim deploy`
if [ ! ${STIM_DEPLOY+x} ]; then echo "Must be run with 'stim deploy'"; exit 1; fi

echo "Deploying Grafana to ${DEPLOY_ENVIRONMENT} in instance ${DEPLOY_INSTANCE} in cluster ${DEPLOY_CLUSTER}"

helm upgrade --install \
  --namespace ${NAMESPACE} \
  --version ${HELM_CHART_VERSION} \
  --set adminPassword=${ADMIN_PASS} \
  --set env.GF_DATABASE_USER="${GF_DATABASE_USER}" \
  --set env.GF_DATABASE_PASSWORD="${GF_DATABASE_PASSWORD}" \
  ${HELM_DEPLOYMENT_NAME} \
  stable/grafana
