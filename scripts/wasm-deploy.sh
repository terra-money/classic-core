#!/bin/sh

BINARY=_build/old/terrad
CONTRACTPATH="contrib/localnet/simulation/misc/cw721_base.wasm"
KEYRING_BACKEND="test"
HOME=mytestnet
CHAIN_ID=test

echo "SETTING UP SMART CONTRACT INTERACTION"

# create two contracts only
for j in $(seq 0 1); do

	echo "key test$j ..."

	# stores contract
	echo "... stores a wasm"
	addr=$($BINARY keys show test$j -a --home $HOME --keyring-backend $KEYRING_BACKEND)
	out=$($BINARY tx wasm store ${CONTRACTPATH} --from test$j --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
	code=$(echo $out | jq -r '.code')
	if [ "$code" != "0" ]; then
		echo "... Could not store NFT binary" >&2
		echo $out >&2
		exit $code
	fi
	sleep 10
	txhash=$(echo $out | jq -r '.txhash')
	id=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[1].attributes[1].value')

	# instantiates contract
	echo "... instantiates contract"
	msg='{"name":"BaseNFT","symbol":"BASE","minter":"'$addr'"}'
	out=$($BINARY tx wasm instantiate $id "$msg" --from test$j --output json --gas auto --gas-adjustment 2.3 --fees 20000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
	code=$(echo $out | jq -r '.code')
	if [ "$code" != "0" ]; then
		echo "... Could not instantiate NFT contract" >&2
		echo $out >&2
		exit $code
	fi
	sleep 10
	txhash=$(echo $out | jq -r '.txhash')
	contract_addr=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[0].attributes[3].value')

	# mints some tokens
	echo "... mints tokens"
	for i in $(seq 0 2); do
		echo "	- token id: "$i
		msg='{"mint":{"token_id":"'$i'","owner":"'$addr'"}}'
		out=$($BINARY tx wasm execute $contract_addr "$msg" --from test$j --output json --gas auto --gas-adjustment 2.3 --fees 20000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
		code=$(echo $out | jq -r '.code')
		if [ "$code" != "0" ]; then
			echo "... Could not mint tokens from contract" $contract_addr >&2
			echo $out >&2
			exit $code
		fi

		sleep 10
	done

	# sends token to other nodes
	echo "... send tokens"
	for i in $(seq 0 2); do
		peer_addr=$($BINARY keys show test$i -a --home $HOME --keyring-backend $KEYRING_BACKEND)
		if [ "$peer_addr" = "$addr" ]; then
			continue
		fi
		msg='{"transfer_nft":{"recipient":"'$peer_addr'","token_id":"'$i'"}}'
		out=$($BINARY tx wasm execute $contract_addr "$msg" --from test$j --output json --gas auto --gas-adjustment 2.3 --fees 20000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
		code=$(echo $out | jq -r '.code')
		if [ "$code" != "0" ]; then
			echo "... Could not transfer NFT id $i from $addr to $peer_addr (contract: $contract_addr)" >&2
			echo $out >&2
			exit $code
		fi

		sleep 10
	done

done
