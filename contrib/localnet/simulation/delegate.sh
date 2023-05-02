#!/bin/sh

# get all validators on the network
VALIDATORS=($(terrad q staking validators --node $(sh $SIMULATION_FOLDER/next_node.sh) -o json | jq -r '.validators[].operator_address'))

# Loop through each node* folder
for operator_address in ${VALIDATORS[@]}
do
    for i in $(seq 1 3); do
        # check balances of test$i to see if it has enough to delegate
        balance=$(terrad q bank balances $(terrad keys show test$i -a --keyring-backend $KEYRING_BACKEND --home $NODE_HOME) --node $(sh $SIMULATION_FOLDER/next_node.sh) -o json | jq -r '.balances | if length > 0 then .[] | select(.denom == "uluna").amount else "0" end')
        if [ $balance -lt 50000000 ]; then
            continue
        fi

        terrad tx staking delegate $operator_address 1000000uluna --chain-id $CHAIN_ID --from test$i --gas auto --gas-adjustment 2.3 --fees 20000000uluna --keyring-backend $KEYRING_BACKEND --home $NODE_HOME --node $(sh $SIMULATION_FOLDER/next_node.sh) -y
        if [ ! $? -eq 0 ]; then
            exit 1
        fi
        
        sleep 10
    done
done