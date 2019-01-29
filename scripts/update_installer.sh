#!/bin/bash

VERSION=${1:-master}

mkdir -p bin
cp install.sh bin
cd bin
sed -i 's/^VERSION\=v.*$/VERSION\='${VERSION}'/g' install.sh
SHA_DARWIN=$(sha256sum stim-darwin-${VERSION}.tar.gz | cut -d' ' -f 1)
sed -i 's/^SHA_DARWIN\=.*$/SHA_DARWIN\='${SHA_DARWIN}'/g' install.sh
SHA_LINUX=$(sha256sum stim-darwin-${VERSION}.tar.gz | cut -d' ' -f 1)
sed -i 's/^SHA_LINUX\=.*$/SHA_LINUX\='${SHA_LINUX}'/g' install.sh
