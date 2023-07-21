#!/bin/sh

BINARY=_build/old/terrad
CONTRACTPATH="scripts/wasm/contracts/old_anc_token.wasm"
KEYRING_BACKEND="test"
HOME=mytestnet
CHAIN_ID=localterra

# upload old contracts
echo "... stores a wasm"
addr=$($BINARY keys show test0 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
addr2=$($BINARY keys show test1 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
out=$($BINARY tx wasm store ${CONTRACTPATH} --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not store binary" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')
id=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[1].attributes[1].value')
echo "CODE = $id"
echo ""

# upload old contracts
echo "... stores a second wasm"
addr=$($BINARY keys show test0 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
addr2=$($BINARY keys show test1 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
out=$($BINARY tx wasm store ${CONTRACTPATH} --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not store binary" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')
# commented out on purpose
# we wanna use the first id to instantiate
# a contract from
# id=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[1].attributes[1].value')
echo "CODE = $id"
echo ""

# instantiates contract
echo "... instantiates contract"
msg=$(jq -n '
{
    "decimals":8,
    "initial_balances":[
        {
            "address":"'$addr'",
            "amount":"1000000000"
        }
    ],
    "name":"Anchor Token",
    "symbol":"ANC"
}')
echo $msg
out=$($BINARY tx wasm instantiate $id "$msg" --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 20000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not instantiate contract" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')

PRE_UPGRADE_CONTRACT_ADDR=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[0].attributes[3].value')
export PRE_UPGRADE_CONTRACT_ADDR
