#!/usr/bin/env bash

set -eo pipefail

# check if there is protoc binary
if ! command -v protoc &> /dev/null
then
  echo "protoc could not be found. Please install http://google.github.io/proto-lens/installing-protoc.html"
  exit
fi

mkdir -p ./tmp-swagger-gen

# Get the path of the cosmos-sdk repo from go/pkg/mod
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
ibc_go_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/ibc-go/v4)
wasm_dir=$(go list -f '{{ .Dir }}' -m github.com/CosmWasm/wasmd)

proto_dirs=$(find ./proto "$cosmos_sdk_dir"/proto "$ibc_go_dir"/proto "$wasm_dir"/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    protoc \
    -I "proto" \
    -I "$cosmos_sdk_dir/third_party/proto" \
    -I "$cosmos_sdk_dir/proto" \
    -I "$ibc_go_dir/proto" \
    -I "$ibc_go_dir/third_party/proto" \
    -I "$wasm_dir/proto" \
    "$query_file" \
    --swagger_out=./tmp-swagger-gen \
    --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
cd ./client/docs
npm install
npm run-script combine
cd ../../

# clean swagger files
rm -rf ./tmp-swagger-gen
