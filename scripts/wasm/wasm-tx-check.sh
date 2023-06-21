#!/bin/bash

# expecting that TXHASH from wasm-deploy.sh will be exported
# querying TXHASH after upgrade to see if it still works

set +e

read -r -a OLD_TXHASH <<< ${TXHASH_STRING:-""}

echo "OLD_TXHASH = ${OLD_TXHASH[@]}"

# loop through OLD_TXHASH
for i in "${OLD_TXHASH[@]}"; do
    echo "querying $i"
    ./_build/new/terrad q tx $i --output json
done