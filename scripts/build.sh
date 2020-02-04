#!/bin/bash

set -eo pipefail

VERSION=${1:-master}
GOOS=${2:-linux}
DOCKER_REPO="premiereglobal/stim"

# Directory to house our binaries
mkdir -p bin

# Build the container
docker build --build-arg VERSION=${VERSION} --build-arg GOOS=${GOOS} -t ${DOCKER_REPO}:${VERSION}-${GOOS} ./

# Extract the binary from the container
docker run --rm --entrypoint "" --name stim-build -v $(pwd)/bin:/stim-bin ${DOCKER_REPO}:${VERSION}-${GOOS} sh -c "cp /usr/bin/stim /stim-bin"

# Zip up the binary
cd bin
tar -cvzf stim-${GOOS}-${VERSION}.tar.gz stim
cd ..

# Build the deploy container
# This command uses the image we just built above as the base image
if [[ ${GOOS} ==  "linux" ]]; then
  docker build --build-arg STIM_IMAGE=${DOCKER_REPO}:${VERSION}-${GOOS} -f Dockerfile.deploy -t ${DOCKER_REPO}:${VERSION}-deploy-${GOOS} ./
fi
