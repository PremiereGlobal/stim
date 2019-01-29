#!/bin/sh

VERSION=${1:-master}
GOOS=${2:-linux}
DOCKER_REPO="readytalk/stim"

# Directory to house our binaries
mkdir -p bin

# Build the binary in Docker
docker build --build-arg VERSION=${VERSION} --build-arg GOOS=${GOOS} -t ${DOCKER_REPO}:${VERSION}-${GOOS} ./

# Run the container in the background in order to extract the binary
docker run --rm --entrypoint "" --name stim-build -d ${DOCKER_REPO}:${VERSION}-${GOOS} sh -c "sleep 120"

docker cp stim-build:/usr/bin/stim bin
docker stop stim-build

# Zip up the binary
cd bin
tar -cvzf stim-${GOOS}-${VERSION}.tar.gz stim

# Get us back to the root
cd ..
