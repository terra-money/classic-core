#!/bin/sh
# during one of simulation, the node that receives transaction died, probably due to having to handle too much.
# this script is to choose the next node for transaction to use

RPC=($(jq -r --arg chain_id "$CHAIN_ID" '.[$chain_id].rpcs[]' $SIMULATION_FOLDER/network/network.json))
TOTAL_NODE=${#RPC[@]}

retry=0
while true; do
    if [ $retry -eq 3 ]; then
        echo "Maximum retry reached, cannot choose next active node..."
        exit 1
    fi

    NODE=${RPC[$((RANDOM % TOTAL_NODE))]}

    # check if chosen node is alive
    curl -s $NODE/status &> /dev/null

    if [ $? -eq 0 ]; then
        echo $NODE
        exit 0
    else
        retry=$((retry + 1))
        sleep 5
    fi
done