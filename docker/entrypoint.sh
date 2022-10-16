#!/bin/sh

# Default to "data".
DATADIR="${DATADIR:-/terra/.terra/data}"
MONIKER="${MONIKER:-docker-node}"
ENABLE_LCD="${ENABLE_LCD:-true}"
MINIMUM_GAS_PRICES=${MINIMUM_GAS_PRICES-0.01133uluna,0.15uusd,0.104938usdr,169.77ukrw,428.571umnt,0.125ueur,0.98ucny,16.37ujpy,0.11ugbp,10.88uinr,0.19ucad,0.14uchf,0.19uaud,0.2usgd,4.62uthb,1.25usek,1.25unok,0.9udkk,2180.0uidr,7.6uphp,1.17uhkd}
SNAPSHOT_NAME="${SNAPSHOT_NAME}"
SNAPSHOT_BASE_URL="${SNAPSHOT_BASE_URL:-https://getsfo.quicksync.io}"
ENABLE_UNSAFE_CORS="${VALIDATOR_ENABLE_UNSAFE_CORS:-false}"
CORS_ALLOWED_ORIGINS="${CORS_ALLOWED_ORIGINS:-[]}"
CORS_ALLOWED_METHODS="${CORS_ALLOWED_METHODS:-["HEAD", "GET", "POST"]}"
CORS_ALLOWED_HEADERS="${CORS_ALLOWED_HEADERS:-["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"]}"

# First sed gets the app.toml moved into place.
# app.toml updates
sed 's/minimum-gas-prices = "0uluna"/minimum-gas-prices = "'"$MINIMUM_GAS_PRICES"'"/g' ~/app.toml > ~/.terra/config/app.toml
sed -i 's/enable-unsafe-cors = false/enable-unsafe-cors = '$ENABLE_UNSAFE_CORS'/g' ~/.terra/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = '$ENABLE_UNSAFE_CORS'/g' ~/.terra/config/app.toml

# Needed to use awk to replace this multiline string.
if [ "$ENABLE_LCD" = true ] ; then
  gawk -i inplace '/^# Enable defines if the API server should be enabled./,/^enable = false/{if (/^enable = false/) print "# Enable defines if the API server should be enabled.\nenable = true"; next} 1' ~/.terra/config/app.toml
fi

# config.toml updates

sed 's/moniker = "moniker"/moniker = "'"$MONIKER"'"/g' ~/config.toml > ~/.terra/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' ~/.terra/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = '"$CORS_ALLOWED_ORIGINS"'/g' ~/.terra/config/config.toml
sed -i 's/cors_allowed_methods = \["HEAD", "GET", "POST", \]/cors_allowed_methods = '"$CORS_ALLOWED_METHODS"'/g' ~/.terra/config/config.toml
sed -i 's/cors_allowed_headers = \["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time", \]/cors_allowed_headers = '"$CORS_ALLOWED_HEADERS"'/g' ~/.terra/config/config.toml

if [ "$CHAINID" = "columbus-5" ] && [[ ! -z "$SNAPSHOT_NAME" ]] ; then 
  # Download the snapshot if data directory is empty.
  res=$(find "$DATADIR" -name "*.db")
  if [ "$res" ]; then
      echo "data directory is NOT empty, skipping quicksync"
  else
      echo "starting snapshot download"
      mkdir -p $DATADIR
      cd $DATADIR
      FILENAME="$SNAPSHOT_NAME"

      # Download
      aria2c -x5 $SNAPSHOT_BASE_URL/$FILENAME
      # Extract
      lz4 -d $FILENAME | tar xf -

      # # cleanup
      rm $FILENAME
  fi
fi

terrad start $TERRAD_STARTUP_PARAMETERS &

#Wait for Terrad to catch up
while true
do
  if ! (( $(echo $(terrad status) | awk -F '"catching_up":|},"ValidatorInfo"' '{print $2}') ));
  then
    break
  fi
  sleep 1
done

if [ ! -z "$VALIDATOR_AUTO_CONFIG" ] && [ "$VALIDATOR_AUTO_CONFIG" = "1" ]; then
  if [ ! -z "$VALIDATOR_KEYNAME" ] && [ ! -z "$VALIDATOR_MNENOMIC" ] && [ ! -z "$VALIDATOR_PASSPHRASE" ] ; then
    terrad keys add $VALIDATOR_KEYNAME --recover > ~/.terra/keys.log 2>&1 << EOF
$VALIDATOR_MNENOMIC
$VALIDATOR_PASSPHRASE
$VALIDATOR_PASSPHRASE
EOF
  fi

  if [ ! -z "$VALIDATOR_AMOUNT" ] && [ ! -z "$MONIKER" ] && [ ! -z "$VALIDATOR_PASSPHRASE" ] && [ ! -z "$VALIDATOR_KEYNAME" ] && [ ! -z "$VALIDATOR_KEYNAME" ] && [ ! -z "$VALIDATOR_COMMISSION_RATE" ] && [ ! -z "$VALIDATOR_COMMISSION_RATE_MAX" ]  && [ ! -z "$VALIDATOR_COMMISSION_RATE_MAX_CHANGE" ]  && [ ! -z "$VALIDATOR_MIN_SELF_DELEGATION" ] ; then
    terrad tx staking create-validator --amount=$VALIDATOR_AMOUNT --pubkey=$(terrad tendermint show-validator) --moniker="$MONIKER" --chain-id=$CHAINID --from=$VALIDATOR_KEYNAME --commission-rate="$VALIDATOR_COMMISSION_RATE" --commission-max-rate="$VALIDATOR_COMMISSION_RATE_MAX" --commission-max-change-rate="$VALIDATOR_COMMISSION_RATE_MAX_CHANGE" --min-self-delegation="$VALIDATOR_MIN_SELF_DELEGATION" --gas=$VALIDATOR_GAS --gas-adjustment=$VALIDATOR_GAS_ADJUSTMENT --fees=$VALIDATOR_FEES > ~/.terra/validator.log 2>&1 << EOF
$VALIDATOR_PASSPHRASE
y
EOF
  fi
fi
wait
