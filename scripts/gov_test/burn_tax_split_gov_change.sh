#!/bin/bash

HOME=mytestnet
ROOT=$(pwd)
DENOM=uluna
KEY="test"
KEY1="test1"
KEY2="test2"
KEYRING="test"
CHAIN_ID="test"

pkill terrad

# underscore so that go tool will not take gocache into account
mkdir -p _build/gocache
export GOMODCACHE=$ROOT/_build/gocache

# install new binary
if ! command -v _build/new/terrad &> /dev/null
then
    GOBIN="$ROOT/_build/new" go install -mod=readonly ./...
fi

# start a node
screen -L -Logfile mytestnet/log-screen.txt -dmS node1 bash scripts/run-node.sh _build/new/terrad $DENOM

sleep 20

SPLIT_RATE=$(./_build/new/terrad q treasury params --output=json | jq '.params.burn_tax_split')
echo "Before split rate: $SPLIT_RATE"

# submit params change proposal
./_build/new/terrad tx gov submit-proposal param-change "$ROOT/scripts/gov_test/burn_tax_split_prop.json" --from $KEY --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME -y

sleep 5

./_build/new/terrad tx gov deposit 1 "20000000${DENOM}" --from $KEY --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME -y

sleep 5

./_build/new/terrad tx gov vote 1 yes --from $KEY --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME -y

sleep 5

./_build/new/terrad tx gov vote 1 yes --from $KEY1 --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME -y

sleep 5

while true; do 
    PROPOSAL_STATUS=$(./_build/new/terrad q gov proposal 1 --output=json | jq ".status" -r)
    echo $PROPOSAL_STATUS
    if [ $PROPOSAL_STATUS = "PROPOSAL_STATUS_PASSED" ]; then
        break
    else
        sleep 10
    fi
done

# check param again
SPLIT_RATE=$(./_build/new/terrad q treasury params --output=json | jq '.params.burn_tax_split')
echo "After split rate: $SPLIT_RATE"