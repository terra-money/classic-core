#!/bin/bash

# $(curl --silent "https://api.github.com/repos/classic-terra/core/releases/latest" | jq -r '.tag_name')
OLD_VERSION=v1.0.5
UPGRADE_HEIGHT=20
HOME=mytestnet
ROOT=$(pwd)
DENOM=uluna
SOFTWARE_UPGRADE_NAME=$(ls -td -- ./app/upgrades/* | head -n 1 | cut -d'/' -f4)

# underscore so that go tool will not take gocache into account
mkdir -p _build/gocache
export GOMODCACHE=$ROOT/_build/gocache

# install old binary
if ! command -v _build/old/terrad &> /dev/null
then
    mkdir -p _build/old
    wget -c "https://github.com/classic-terra/core/archive/refs/tags/${OLD_VERSION}.zip" -O _build/${OLD_VERSION}.zip
    unzip _build/${OLD_VERSION}.zip -d _build
    cd ./_build/core-${OLD_VERSION:1}
    GOBIN="$ROOT/_build/old" go install -mod=readonly ./...
    cd ../..
fi

# reinstall old binary
if [ $# -eq 1 ] && [ $1 == "--reinstall-old" ]; then
    cd ./_build/core-${OLD_VERSION:1}
    GOBIN="$ROOT/_build/old" go install -mod=readonly ./...
    cd ../..
fi

# install new binary
if ! command -v _build/new/terrad &> /dev/null
then
    GOBIN="$ROOT/_build/new" go install -mod=readonly ./...
fi

# start old node
screen -L -Logfile mytestnet/log-screen.txt -dmS node1 bash scripts/run-node.sh _build/old/terrad $DENOM

sleep 20

./_build/old/terrad tx gov submit-proposal software-upgrade "$SOFTWARE_UPGRADE_NAME" --upgrade-height $UPGRADE_HEIGHT --upgrade-info "temp" --title "upgrade" --description "upgrade"  --from test1 --keyring-backend test --chain-id test --home $HOME -y

sleep 5

./_build/old/terrad tx gov deposit 1 "20000000${DENOM}" --from test1 --keyring-backend test --chain-id test --home $HOME -y

sleep 5

./_build/old/terrad tx gov vote 1 yes --from test --keyring-backend test --chain-id test --home $HOME -y

sleep 5

./_build/old/terrad tx gov vote 1 yes --from test1 --keyring-backend test --chain-id test --home $HOME -y

sleep 5

# determine block_height to halt
while true; do 
    BLOCK_HEIGHT=$(./_build/old/terrad status | jq '.SyncInfo.latest_block_height' -r)
    if [ $BLOCK_HEIGHT = "$UPGRADE_HEIGHT" ]; then
        # assuming running only 1 terrad
        echo "BLOCK HEIGHT = $UPGRADE_HEIGHT REACHED, KILLING OLD ONE"
        pkill terrad
        break
    else
        ./_build/old/terrad q gov proposal 1 --output=json | jq ".status"
        echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
        sleep 10
    fi
done

sleep 5

./_build/new/terrad start --home $HOME