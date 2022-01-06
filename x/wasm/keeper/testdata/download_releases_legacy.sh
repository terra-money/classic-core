# https://github.com/terra-money/core/raw/main/x/wasm/keeper/testdata/bindings_tester.wasm

#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

if [ $# -ne 1 ]; then
  echo "Usage: ./download_legacy_contracts.sh REV"
  exit 1
fi

rev="$1"

for contract in burner hackatom reflect staking bindings_tester maker; do
  url="https://github.com/terra-money/core/raw/$rev/x/wasm/keeper/testdata/${contract}.wasm"
  echo "Downloading $url ..."
  wget -O "${contract}_legacy.wasm" "$url"
done

# create the zip variant
gzip -k hackatom_legacy.wasm
mv hackatom_legacy.wasm.gz hackatom_legacy.wasm.gzip

rm -f version_legacy.txt
echo "$rev" > version_legacy.txt
