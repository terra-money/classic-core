#!/bin/bash

VERSION="${1:-v0.5.21}"

pushd .. 

#git checkout $VERSION
docker build -t terrarebels/classic-core:$VERSION .
#git checkout -

popd

#docker build --build-arg version=$VERSION --build-arg chainid=columbus-5 -t terrarebels/classic-core-node:$VERSION .
docker build --build-arg version=$VERSION --build-arg chainid=rebel-1 -t terrarebels/classic-core-node:$VERSION-testnet .