#!/bin/bash

set -eo pipefail

VERSION=${1:-v0.0.0-local}
GOOS=${2:-linux}
DOCKER_REPO="premiereglobal/stim"

if [[ ${VERSION:0:1} != "v" ]]; then
  echo "VERSION must start with a v (ie v0.0.0-branch) and is currently ${VERSION}"
  exit 1
fi

VA="${VERSION//[^.]}"
VC="${#VA}"

if [[ ${VC} != "2" ]]; then
  echo "VERSION must be in syntax 'v1.1.1' or 'v1.1.1-alpha' only 3 version numbers sperated by a '.' and a build/prerelease string starting with a '-' at the end, currently ${VERSION} ${VC}"
  exit 1
fi

# Directory to house our binaries
mkdir -p bin

echo "Building version:\"${VERSION}\""

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
