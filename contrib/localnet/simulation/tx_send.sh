#!/bin/sh

# randomly pick two addresses in keyring and send 100000000uluna
addresses=($(terrad keys list -n --keyring-backend $KEYRING_BACKEND --home $NODE_HOME))
length=${#addresses[@]}

success=0
while true; do
    if [ $success -eq 8 ]; then
        break
    fi

    addr_name_1=${addresses[$RANDOM % length]}
    addr_name_2=${addresses[$RANDOM % length]}
    if [ "$addr_name_1" == "$addr_name_2" ]; then
        continue
    fi

    addr1=$(terrad keys show $addr_name_1 -a --keyring-backend $KEYRING_BACKEND --home $NODE_HOME)
    addr2=$(terrad keys show $addr_name_2 -a --keyring-backend $KEYRING_BACKEND --home $NODE_HOME)

    # check balances of addr1 and addr2
    balance1=$(terrad q bank balances $addr1 --node $(sh $SIMULATION_FOLDER/next_node.sh) -o json | jq -r '.balances | if length > 0 then .[] | select(.denom == "uluna").amount else "0" end')
    if [ $balance1 -lt 500000000 ]; then
        continue
    fi

    terrad tx bank send $addr1 $addr2 1000000uluna --chain-id $CHAIN_ID --home $NODE_HOME --gas auto --gas-adjustment 2.3 --fees 20000000uluna --keyring-backend $KEYRING_BACKEND --node $(sh $SIMULATION_FOLDER/next_node.sh) -y
    if [ $? -eq 0 ]; then
        success=$((success+1))
    fi
    sleep 10
done