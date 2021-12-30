#!/bin/sh

# Default to "data".
DATADIR="${DATADIR:-/root/.terra/data}"
MONIKER="${MONIKER:-docker-node}"
ENABLE_LCD="${ENABLE_LCD:-true}"
MINIMUM_GAS_PRICES=${MINIMUM_GAS_PRICES-0.01133uluna,0.15uusd,0.104938usdr,169.77ukrw,428.571umnt,0.125ueur,0.98ucny,16.37ujpy,0.11ugbp,10.88uinr,0.19ucad,0.14uchf,0.19uaud,0.2usgd,4.62uthb,1.25usek,1.25unok,0.9udkk,2180.0uidr,7.6uphp,1.17uhkd}
SNAPSHOT_NAME="${SNAPSHOT_NAME}"
SNAPSHOT_BASE_URL="${SNAPSHOT_BASE_URL:-https://getsfo.quicksync.io}"

# First sed gets the app.toml moved into place.
# app.toml updates
sed 's/minimum-gas-prices = "0uluna"/minimum-gas-prices = "'"$MINIMUM_GAS_PRICES"'"/g' ~/app.toml > ~/.terra/config/app.toml

# Needed to use awk to replace this multiline string.
if [ "$ENABLE_LCD" = true ] ; then
  gawk -i inplace '/^# Enable defines if the API server should be enabled./,/^enable = false/{if (/^enable = false/) print "# Enable defines if the API server should be enabled.\nenable = true"; next} 1' ~/.terra/config/app.toml
fi

# config.toml updates

sed 's/moniker = "moniker"/moniker = "'"$MONIKER"'"/g' ~/config.toml > ~/.terra/config/config.toml
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' ~/.terra/config/config.toml

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

exec "$@" --db_dir $DATADIR