#!/bin/bash

#
# start network
TESTNET_NVAL=6 TESTNET_VOTING_PERIOD=40s make clean localnet-start

#
# convenience
NODE1_HOME=node1/terrad
TERRAD="docker exec terradnode1 ./terrad"
echo "{\"title\":\"prop1\",\"description\":\"prop1\",\"changes\":[{\"subspace\":\"treasury\",\"key\":\"MinInitialDepositRatio\",\"value\":\"0.1\"}]}" > build/node1/terrad/prop.json
echo "{\"title\":\"prop2\",\"description\":\"prop2\",\"changes\":[{\"subspace\":\"treasury\",\"key\":\"MinInitialDepositRatio\",\"value\":\"0.0\"}]}" > build/node1/terrad/prop2.json
echo "{\"title\":\"prop2\",\"description\":\"prop2\",\"changes\":[{\"subspace\":\"treasury\",\"key\":\"MinInitialDepositRatio\",\"value\":\"0.0\"}],\"deposit\":\"100000uluna\"}" > build/node1/terrad/prop3.json
echo "{\"title\":\"prop2\",\"description\":\"prop2\",\"changes\":[{\"subspace\":\"treasury\",\"key\":\"MinInitialDepositRatio\",\"value\":\"0.0\"}],\"deposit\":\"1000000uluna\"}" > build/node1/terrad/prop4.json

#
# sleep to wait for localnet to come up
sleep 10

#
# 100 block from now
STATUS_INFO=($($TERRAD status --home $NODE1_HOME | jq -r '.NodeInfo.network,.SyncInfo.latest_block_height'))
CHAIN_ID=${STATUS_INFO[0]}

#
# initial prop should pass and increases min initial deposit ratio to 0.1
STATUS=$($TERRAD tx gov submit-proposal param-change $NODE1_HOME/prop.json --keyring-backend test --yes --from node1 --broadcast-mode block --chain-id $CHAIN_ID --home $NODE1_HOME -o json | jq -r .code)
if [ $STATUS -ne 0 ]; then
	echo "Testcase1: failed to submit prop with no initial deposit"
	exit -1
else
	echo "Testcase1: succeeded to submit prop with no initial deposit"
fi

#
# provide deposit
sleep 5
echo "yproviding deposit"
$TERRAD tx gov deposit 1 10000000uluna --from node1 --keyring-backend test --chain-id $CHAIN_ID --home "node1/terrad" -y

#
# loop from 0 to 5 to vote yes
for i in {0..5}; do
    # check if docker for node i is running
    if [[ $(docker ps -a | grep terradnode$i | wc -l) -eq 1 ]]; then
		sleep 5
        $TERRAD tx gov vote 1 yes --from node$i --keyring-backend test --chain-id $CHAIN_ID --home "node$i/terrad" -y
    fi
done

#
# wait for proposal to pass
for i in {0..6}; do
    STATUS=$($TERRAD q gov proposal 1 --chain-id $CHAIN_ID --home $NODE1_HOME -o json | jq -r .status)
    echo $STATUS
    if [ "$STATUS" = "PROPOSAL_STATUS_PASSED" ]; then
		break
	fi
    sleep 10
done

sleep 10

#
# second prop should fail without deposit ratio to 0.1
STATUS=$($TERRAD tx gov submit-proposal param-change $NODE1_HOME/prop2.json --keyring-backend test --yes --from node1 --broadcast-mode block --chain-id $CHAIN_ID --home $NODE1_HOME -o json | jq -r .code)
if [ $STATUS -ne 0 ]; then
	echo "Testcase2: failed to submit prop with no initial deposit... as expected"
else
	echo "Testcase2: succeeded to submit prop with no initial deposit... unexpectedly"
	exit -1
fi

sleep 10

#
# retry - this time with 100000 uluna - insufficient deposit
STATUS=$($TERRAD tx gov submit-proposal param-change $NODE1_HOME/prop3.json --keyring-backend test --yes --from node1 --broadcast-mode block --chain-id $CHAIN_ID --home $NODE1_HOME -o json | jq -r .code)
if [ $STATUS -ne 0 ]; then
	echo "Testcase3: failed to submit prop with insufficient initial deposit... as expected"
else
	echo "Testcase3: succeeded to submit prop with insufficient initial deposit... unexpectedly"
	exit -1
fi

#
# retry - this time with 1000000 uluna - sufficient deposit
STATUS=$($TERRAD tx gov submit-proposal param-change $NODE1_HOME/prop4.json --keyring-backend test --yes --from node1 --broadcast-mode block --chain-id $CHAIN_ID --home $NODE1_HOME -o json | jq -r .code)
if [ $STATUS -ne 0 ]; then
	echo "Testcase4: failed to submit prop with initial deposit... unexpectedly"
	exit -1
else
	echo "Testcase4: succeeded to submit prop with initial deposit... as expected"
fi

echo "all good"

exit 0
