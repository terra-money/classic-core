#!/bin/bash

# should make this auto fetch upgrade name from app upgrades once many upgrades have been done
# this command will fetch the latest upgrade folder by time modified. So, there could be case that v1 is later modified than v2 which will cause error
SOFTWARE_UPGRADE_NAME=$(ls -td -- ./app/upgrades/*/ | head -n 1 | xargs basename)
NODE1_HOME=node1/terrad
BINARY_OLD="docker exec terradnode1 ./old/terrad"
TESTNET_NVAL=${1:-7}

# sleep to wait for localnet to come up
sleep 10

# 100 block from now
STATUS_INFO=($($BINARY_OLD status --home $NODE1_HOME | jq -r '.NodeInfo.network,.SyncInfo.latest_block_height'))
echo $STATUS_INFO
CHAIN_ID=${STATUS_INFO[0]}
UPGRADE_HEIGHT=$((STATUS_INFO[1] + 20))
echo $UPGRADE_HEIGHT

$BINARY_OLD tx gov submit-proposal software-upgrade "$SOFTWARE_UPGRADE_NAME" --upgrade-height $UPGRADE_HEIGHT --upgrade-info "temp" --title "upgrade" --description "upgrade"  --from node1 --keyring-backend test --chain-id $CHAIN_ID --home $NODE1_HOME -y

sleep 5

$BINARY_OLD tx gov deposit 1 "20000000uluna" --from node1 --keyring-backend test --chain-id $CHAIN_ID --home $NODE1_HOME -y

sleep 5

# loop from 0 to TESTNET_NVAL
for (( i=0; i<$TESTNET_NVAL; i++ )); do
    # check if docker for node i is running
    if [[ $(docker ps -a | grep terradnode$i | wc -l) -eq 1 ]]; then
        $BINARY_OLD tx gov vote 1 yes --from node$i --keyring-backend test --chain-id $CHAIN_ID --home "node$i/terrad" -y
        sleep 5
    fi
done

# keep track of block_height
NIL_BLOCK=0
LAST_BLOCK=0
SAME_BLOCK=0
while true; do 
    BLOCK_HEIGHT=$($BINARY_OLD status --home $NODE1_HOME | jq '.SyncInfo.latest_block_height' -r)
    # if BLOCK_HEIGHT is empty
    if [[ -z $BLOCK_HEIGHT ]]; then
        # if 5 nil blocks in a row, exit
        if [[ $NIL_BLOCK -ge 5 ]]; then
            echo "ERROR: 5 nil blocks in a row"
            break
        fi
        NIL_BLOCK=$((NIL_BLOCK + 1))
    fi

    # if block height is not nil
    # if block height is same as last block height
    if [[ $BLOCK_HEIGHT -eq $LAST_BLOCK ]]; then
        # if 5 same blocks in a row, exit
        if [[ $SAME_BLOCK -ge 5 ]]; then
            echo "ERROR: 5 same blocks in a row"
            break
        fi
        SAME_BLOCK=$((SAME_BLOCK + 1))
    else
        # update LAST_BLOCK and reset SAME_BLOCK
        LAST_BLOCK=$BLOCK_HEIGHT
        SAME_BLOCK=0
    fi

    if [[ $BLOCK_HEIGHT -ge $UPGRADE_HEIGHT ]]; then
        # assuming running only 1 terrad
        echo "UPGRADE REACHED, CONTINUING NEW CHAIN"
        break
    else
        $BINARY_OLD q gov proposal 1 --output=json --home $NODE1_HOME | jq ".status"
        echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
        sleep 10
    fi
done

if [[ $SAME_BLOCK -ge 5 ]]; then
    docker logs terradnode0
    exit 1
fi

sleep 40

# check all nodes are online after upgrade
for (( i=0; i<$TESTNET_NVAL; i++ )); do
    if [[ $(docker ps -a | grep terradnode$i | wc -l) -eq 1 ]]; then
        docker exec terradnode$i ./terrad status --home "node$i/terrad"
        if [[ "${PIPESTATUS[0]}" != "0" ]]; then
            echo "node$i is not online"
            docker logs terradnode$i
            exit 1
        fi
    else
        echo "terradnode$i is not running"
        docker logs terradnode$i
        exit 1
    fi
done