#!/bin/sh

set -e

VERSION=v0.0.2
SHA_LINUX=b39c570c3a6e2e5a114fcf89f061ff9fb912527242f2a80f2b9248a567ebf2b6
SHA_DARWIN=0dbae3a2c61d93f69c4f9274ee5227d1969d19a4f9672d07d39fe738e280e085
CACHE_DIR=${HOME}/.stim/cache
BIN_DIR=${HOME}/.stim/bin

# Extracted file already exists
# Assume it's fine (for now)
# TODO: Verify file somehow?
if [ -f ${BIN_DIR}"/stim-${VERSION}" ]; then
  exit 0
fi

# Verify Signature of file
verify() {

  FILE=${1}
  SHA=${2}

  which shasum 2>&1 > /dev/null
  if [ $? -eq 0 ]; then
    if [ "$(shasum -a 256 ${FILE} | cut -d' ' -f 1)" != "${SHA}" ]; then
      echo 1
    fi
  fi

  which sha256sum 2>&1 > /dev/null
  if [ $? -eq 0 ]; then
    if [ "$(sha256sum ${FILE} | cut -d' ' -f 1)" != "${SHA}" ]; then
      echo 1
    fi
  fi

  echo 0
}

# Download binary
download() {

  which wget 2>&1 > /dev/null
  if [ $? -eq 0 ]; then
    wget --quiet -O ${ARCHIVE} https://github.com/PremiereGlobal/stim/releases/download/${VERSION}/${ARCHIVE}
    return 0
  fi

  which curl 2>&1 > /dev/null
  if [ $? -eq 0 ]; then
    curl -L -s -o ${ARCHIVE} https://github.com/PremiereGlobal/stim/releases/download/${VERSION}/${ARCHIVE}
    return 0
  fi

  >&2 echo "'wget' or 'curl' not found, cannot download binary"
  exit 1

}

# Change working directory to
mkdir -p ${BIN_DIR}
mkdir -p ${CACHE_DIR}
cd ${CACHE_DIR}

# Determine OS
if [ "${OSTYPE}" = "linux-gnu" -o "$(uname)" = "Linux" ]; then
  OS=linux
  SHA=${SHA_LINUX}
elif [ "$(uname)" = "Darwin" ]; then
  OS=darwin
  SHA=${SHA_DARWIN}
else
  >&2 echo "Could not detect OS - failing"
  exit 1
fi

ARCHIVE=stim-${OS}-${VERSION}.zip

if [ -f ${ARCHIVE} ]; then
  if [ "$(verify ${ARCHIVE} ${SHA})" -eq 0 ]; then
    # Existing valid archive found in cache, use it
    unzip -q -p ${ARCHIVE} > ${BIN_DIR}"/stim-${VERSION}"
    chmod +x ${BIN_DIR}"/stim-${VERSION}"
    exit 0
  fi
fi

# We don't have a valid archive in cache, download it
download

# Verify the downloaded file is valid
if [ $(verify ${ARCHIVE} ${SHA}) -eq 0 ]; then
  unzip -q -p ${ARCHIVE} > ${BIN_DIR}"/stim-${VERSION}"
  chmod +x ${BIN_DIR}"/stim-${VERSION}"
else
  >&2 echo "Signature of downloaded file '"${ARCHIVE}"' is invalid"
fi

exit 0
