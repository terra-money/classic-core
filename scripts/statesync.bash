#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -uxe

# set environment variables
#export GOPATH=~/go
#export PATH=$PATH:~/go/bin
export NODE=65.21.202.37:2161
export APPNAME=TERRAD


# Install Gaia
# go install ./...


# MAKE HOME FOLDER AND GET GENESIS
gaiad init test 
wget -O ~/.gaia/config/genesis.json https://cloudflare-ipfs.com/ipfs/Qmc54DreioPpPDUdJW6bBTYUKepmcPsscfqsfFcFmTaVig

INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s $NODE/height | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT-$INTERVAL))
NODE_ID=$(curl -s "$NODE/block?height=$BLOCK_HEIGHT" | jq -r .result.node_info.id)
TRUST_HASH=$(curl -s "$NODE/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"
echo "NODE ID: $NODE_ID"


# export state sync ars
export $($APPNAME)_STATESYNC_ENABLE=true
export $($APPNAME)_P2P_MAX_NUM_OUTBOUND_PEERS=200
export $($APPNAME)_STATESYNC_RPC_SERVERS="$NODE,https://columbus-5.technofractal.com:443"
export $($APPNAME)_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export $($APPNAME)_STATESYNC_TRUST_HASH=$TRUST_HASH
export $($APPNAME)_P2P_PERSISTENT_PEERS="$NODE_ID@$NODE"

terrad start
