#!/bin/sh
set -eu -o pipefail

# Exit if not running with `stim deploy`
if [ ! ${STIM_DEPLOY+x} ]; then echo "Must be run with 'stim deploy'"; exit 1; fi

# This script is meant to be run from a docker container
helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade \
  --debug \
  --install \
  --namespace ${NAMESPACE} \
  --set image.tag=${IMAGE_TAG} \
  nginx-test bitnami/nginx
