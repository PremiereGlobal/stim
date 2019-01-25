#!/bin/sh

VERSION=v0.0.2
CACHE_DIR=${HOME}/.stim/cache

mkdir -p ${CACHE_DIR}
cd ${CACHE_DIR}

if [[ -f stim-${VERSION} && "$(./stim-${VERSION} version)" == "stim/${VERSION}" ]]; then
  cp stim-${VERSION} stim
  exit 0
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
  wget --quiet -O stim.zip https://github.com/ReadyTalk/stim/releases/download/${VERSION}/stim-linux-${VERSION}.zip
elif [[ "$OSTYPE" == "darwin"* ]]; then
  curl -L -s -o stim.zip https://github.com/ReadyTalk/stim/releases/download/${VERSION}/stim-darwin-${VERSION}.zip
else
  echo "Could not detect OS - failing"
  exit 1
fi

unzip -o stim.zip
rm stim.zip
cp stim stim-${VERSION}

exit 0
