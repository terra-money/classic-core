#/bin/sh

# KEY MANAGEMENT
KEYRING="test"

# Function updates the config based on a jq argument as a string
update_test_genesis () {
    # EX: update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat ~/.terra/config/genesis.json | jq --arg DENOM "$2" "$1" > ~/.terra/config/tmp_genesis.json && mv ~/.terra/config/tmp_genesis.json ~/.terra/config/genesis.json
}

# add keys, add balances
for i in $(seq 0 3); do
    key=$(jq ".keys[$i] | tostring" /keys.json )
    keyname=$(echo $key | jq -r 'fromjson | ."keyring-keyname"')
    mnemonic=$(echo $key | jq -r 'fromjson | .mnemonic')
    # Add new account
    echo $mnemonic | terrad keys add $keyname --keyring-backend $KEYRING --recover --home ~/.terra
    # Add initial balances
    terrad add-genesis-account $keyname "1000000000000uluna" --keyring-backend $KEYRING --home ~/.terra
done

# Sign genesis transaction
terrad gentx test "1000000uluna" --keyring-backend $KEYRING --chain-id $CHAINID --home ~/.terra

update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"] = "50s"'
update_test_genesis '.app_state["mint"]["params"]["mint_denom"]=$DENOM' uluna
update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": $DENOM,"amount": "1000000"}]' uluna
update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": $DENOM,"amount": "1000"}' uluna
update_test_genesis '.app_state["staking"]["params"]["bond_denom"]=$DENOM' uluna

# Collect genesis tx
terrad collect-gentxs --home ~/.terra

# Run this to ensure everything worked and that the genesis file is setup correctly
terrad validate-genesis --home ~/.terra