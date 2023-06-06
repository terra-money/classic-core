#!/bin/bash

rm -rf mytestnet
pkill terrad

BINARY=$1
# check DENOM is set. If not, set to uluna
DENOM=${2:-uluna}

COMMISSION_RATE=0.01
COMMISSION_MAX_RATE=0.02

SED_BINARY=sed
# check if this is OS X
if [[ "$OSTYPE" == "darwin"* ]]; then
    # check if gsed is installed
    if ! command -v gsed &> /dev/null
    then
        echo "gsed could not be found. Please install it with 'brew install gnu-sed'"
        exit
    else
        SED_BINARY=gsed
    fi
fi

# check BINARY is set. If not, build terrad and set BINARY
if [ -z "$BINARY" ]; then
    make build
    BINARY=build/terrad
fi

HOME=mytestnet
CHAIN_ID="localterra"
KEYRING="test"
KEY="test0"
KEY1="test1"
KEY2="test2"

# Function updates the config based on a jq argument as a string
update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/config/genesis.json | jq "$1" > $HOME/config/tmp_genesis.json && mv $HOME/config/tmp_genesis.json $HOME/config/genesis.json
}

$BINARY init --chain-id $CHAIN_ID moniker --home $HOME

$BINARY keys add $KEY --keyring-backend $KEYRING --home $HOME
$BINARY keys add $KEY1 --keyring-backend $KEYRING --home $HOME
$BINARY keys add $KEY2 --keyring-backend $KEYRING --home $HOME

# Allocate genesis accounts (cosmos formatted addresses)
$BINARY add-genesis-account $KEY "1000000000000${DENOM}" --keyring-backend $KEYRING --home $HOME
$BINARY add-genesis-account $KEY1 "1000000000000${DENOM}" --keyring-backend $KEYRING --home $HOME
$BINARY add-genesis-account $KEY2 "1000000000000${DENOM}" --keyring-backend $KEYRING --home $HOME

update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"]="50s"'
update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="'$DENOM'"'
update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom":"'$DENOM'","amount": "1000000"}]'
update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom":"'$DENOM'","amount":"1000"}'
update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="'$DENOM'"'

# enable rest server and swagger
$SED_BINARY -i '0,/enable = false/s//enable = true/' $HOME/config/app.toml
$SED_BINARY -i 's/swagger = false/swagger = true/' $HOME/config/app.toml

# Sign genesis transaction
$BINARY gentx $KEY "1000000${DENOM}" --commission-rate=$COMMISSION_RATE --commission-max-rate=$COMMISSION_MAX_RATE --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME

# Collect genesis tx
$BINARY collect-gentxs --home $HOME

# Run this to ensure everything worked and that the genesis file is setup correctly
$BINARY validate-genesis --home $HOME

$BINARY start --home $HOME