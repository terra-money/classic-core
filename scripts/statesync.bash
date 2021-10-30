#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -uxe

# set environment variables
#export GOPATH=~/go
#export PATH=$PATH:~/go/bin
export NODE=65.21.202.37:2161
export ENVNAME=TERRAD
export APPNAME=terrad
export $(echo $ENVNAME)$(echo "_HOME")=/plotting/terradata
export GENESIS=QmZAMcdu85Qr8saFuNpL9VaxVqqLGWNAs72RVFhchL9jWs

# Install Gaia
# go install ./...


# MAKE HOME FOLDER AND GET GENESIS
$APPNAME init test 
wget -O $TERRAD_HOME/config/genesis.json https://ipfs.io/ipfs/$GENESIS

INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s $NODE/height | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT-$INTERVAL))
NODE_ID=$(curl -s "$NODE/status" | jq -r .result.node_info.id)
TRUST_HASH=$(curl -s "$NODE/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"
echo "NODE ID: $NODE_ID"
echo "HOME: $HOMEDIR"

# export state sync ars
export $(echo $ENVNAME)_STATESYNC_ENABLE=true
export $(echo $ENVNAME)_P2P_MAX_NUM_OUTBOUND_PEERS=200
export $(echo $ENVNAME)_STATESYNC_RPC_SERVERS="$NODE,https://columbus-5.technofractal.com:443"
export $(echo $ENVNAME)_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export $(echo $ENVNAME)_STATESYNC_TRUST_HASH=$TRUST_HASH
export $(echo $ENVNAME)_P2P_PERSISTENT_PEERS="$NODE_ID@$NODE"

terrad start
