#! /bin/bash

BINARY=_build/new/terrad
CONTRACTPATH="scripts/wasm/contracts/new_anc_token.wasm"
RECEIVERPATH="scripts/wasm/contracts/cw20_receiver_template.wasm"
KEYRING_BACKEND="test"
HOME=mytestnet
CHAIN_ID=localterra

if [ -z "$PRE_UPGRADE_CONTRACT_ADDR" ]; then
	echo "PRE_UPGRADE_CONTRACT_ADDR is empty"
	exit 1
fi

### DEBUG ###
#contract_addr="terra18vd8fpwxzck93qlwghaj6arh4p7c5n896xzem5"
### DEBUG ###
addr2=$($BINARY keys show test1 -a --home $HOME --keyring-backend $KEYRING_BACKEND)

echo "STORE DUMMY RECEIVER"
res=$($BINARY tx wasm store $RECEIVERPATH --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --broadcast-mode block --keyring-backend $KEYRING_BACKEND -y)
tx=$(echo $res | jq -r ."txhash")
res=$($BINARY q tx --output json ${tx})
code=$(echo "$res" | jq -r ."code")
if [ "$code" != "0" ]; then
	echo "store contract failed"
	exit -1
fi
code_id=$(echo "$res" | jq -r '.logs[0].events[] | select(.type == "store_code") | .attributes[] | select(.key == "code_id") | .value')
echo "CODE = ${code_id}"
echo ""

echo "INSTANTIATE DUMMY RECEIVER"
res=$($BINARY tx wasm instantiate $code_id '{}' --label "contract_${code_id}" --no-admin --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --broadcast-mode block --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $res | jq -r ."code")
tx=$(echo $res | jq -r ."txhash")
if [ "$code" != 0 ]; then
	echo "instantiate contract failed"
	exit -1
fi
receiver=$($BINARY q tx --output json ${tx} | jq -r '.logs[0].events[] | select(.type == "instantiate") | .attributes[] | select(.key == "_contract_address") | .value')
echo "ADDRESS = ${receiver}"
echo ""

echo "TRANSFER P2P - BEFORE MIGRATION"
msg=$(jq -n '
{
    "transfer": {
		"amount": "100000",
		"recipient": "'$addr2'"
	}
}')
echo $msg
res=$($BINARY tx wasm execute "$PRE_UPGRADE_CONTRACT_ADDR" "$msg" --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
tx=$(echo $res | jq -r ."txhash")
code=$(echo $res | jq -r ."code")
if [ "$code" != "0" ]; then
	echo "transfer message failed"
	exit -1
fi

echo $res

sleep 5

echo "SEND - BEFORE MIGRATION"
msg=$(jq -n '
{
    "send": {
		"amount": "1",
		"contract": "'$receiver'",
		"msg": "eyJ0ZXJtIjp7ImFtb3VudCI6IjEwMDAwMCJ9fQ=="
	}
}')
echo $msg
res=$($BINARY tx wasm execute "$PRE_UPGRADE_CONTRACT_ADDR" "$msg" --from test0 --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
tx=$(echo $res | jq -r ."txhash")
code=$(echo $res | jq -r ."code")
if [ "$code" != "0" ]; then
	echo "transfer message failed"
	exit -1
fi
echo $res