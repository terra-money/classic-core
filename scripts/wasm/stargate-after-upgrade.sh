#!/bin/sh

set +e

BINARY=_build/new/terrad
STARGATE_TESTER="wasmbinding/testdata/stargate_tester.wasm"
KEYRING_BACKEND="test"
HOME=mytestnet
CHAIN_ID=localterra

# store stargate-tester
echo "... stores a wasm"
addr=$($BINARY keys show test0 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
out=$($BINARY tx wasm store ${STARGATE_TESTER} --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not store contract" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')
echo "$txhash"
code_id=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[1].attributes[1].value')
echo "code_id $code_id"

# instantiate stargate-tester
echo "... instantiates contract"
msg='{}'
out=$($BINARY tx wasm instantiate $code_id "$msg" --from test0 --output json --gas auto --gas-adjustment 2.3 --label "stargate-tester" --fees 20000000uluna --no-admin --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not instantiate contract" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')
echo "$txhash"
contract_addr=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[0].attributes[0].value')
echo "contract_addr $contract_addr"

# call stargate-tester
echo "... query tax rate"
msg='{"tax_rate":{}}'
out=$($BINARY query wasm contract-state smart $contract_addr $msg --output json)
echo $out