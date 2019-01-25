#!/bin/sh

set -e

VERSION=v0.0.2
SHA_LINUX=b39c570c3a6e2e5a114fcf89f061ff9fb912527242f2a80f2b9248a567ebf2b6
SHA_DARWIN=0dbae3a2c61d93f69c4f9274ee5227d1969d19a4f9672d07d39fe738e280e085
CACHE_DIR=${HOME}/.stim/cache

# Verify Signature of file
verify() {

  FILE=${1}
  SHA=${2}

  if [[ "$(shasum -a 256 ${FILE} | cut -d' ' -f 1)" != "${SHA}" ]]; then
    echo 1
  fi

  echo 0
}

# Change working directory to
mkdir -p ${CACHE_DIR}
cd ${CACHE_DIR}

# Determine OS
if [[ "${OSTYPE}" == "linux-gnu" ]]; then
  OS=linux
  SHA=${SHA_LINUX}
elif [[ "${OSTYPE}" == "darwin"* ]]; then
  OS=darwin
  SHA=${SHA_DARWIN}
else
  echo "Could not detect OS - failing"
  exit 1
fi

ARCHIVE=stim-${OS}-${VERSION}.zip

if [[ -f ${ARCHIVE} && $(verify ${ARCHIVE} ${SHA}) == 0 ]]; then
  # Existing valid archive found in cache, use it
  unzip -q -o ${ARCHIVE}
  exit 0
fi

# We don't have a valid archive in cache, download it
if [[ "${OSTYPE}" == "linux-gnu" ]]; then
  wget --quiet -O ${ARCHIVE} https://github.com/ReadyTalk/stim/releases/download/${VERSION}/${ARCHIVE}
elif [[ "${OSTYPE}" == "darwin"* ]]; then
  ARCHIVE=stim-darwin-${VERSION}.zip
  curl -L -s -o ${ARCHIVE} https://github.com/ReadyTalk/stim/releases/download/${VERSION}/${ARCHIVE}
fi

# Verify the downloaded file is valid
if [[ $(verify ${ARCHIVE} ${SHA}) == 0 ]]; then
  unzip -q -o ${ARCHIVE}
else
  >&2 echo "Signature of downloaded file '"${ARCHIVE}"' is invalid"
fi

exit 0
