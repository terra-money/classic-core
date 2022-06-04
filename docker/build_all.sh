#!/bin/bash

VERSION="public-node"

pushd .. 

git checkout $VERSION
docker build -t toban/classic-core:$VERSION .
git checkout -

popd

docker build --build-arg version=$VERSION --build-arg chainid=columbus-5 -t toban/classic-core-node:$VERSION .
docker build --build-arg version=$VERSION --build-arg chainid=bombay-12 -t toban/classic-core-node:$VERSION-testnet .