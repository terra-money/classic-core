#!/bin/bash

VERSION="${1:-v0.5.11-oracle}"

pushd .. 

git checkout $VERSION
docker build -t terramoney/core:$VERSION .
git checkout -

popd

docker build --build-arg version=$VERSION --build-arg chainid=columbus-5 -t terramoney/core-node:$VERSION .
docker build --build-arg version=$VERSION --build-arg chainid=bombay-12 -t terramoney/core-node:$VERSION-testnet .