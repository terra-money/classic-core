#!/bin/sh

set +e

BINARY=_build/new/terrad
HACK_ATOM="./custom/auth/ante/testdata/hackatom.wasm"
KEYRING_BACKEND="test"
HOME=mytestnet
CHAIN_ID=localterra
KEY="test0"
KEY1="test1"
KEY1="test2"
DENOM=uluna

# store stargate-tester
echo "... stores a wasm"
addr=$($BINARY keys show $KEY -a --home $HOME --keyring-backend $KEYRING_BACKEND)
addr1=$($BINARY keys show $KEY1 -a --home $HOME --keyring-backend $KEYRING_BACKEND)
addr2=$($BINARY keys show $KEY1 -a --home $HOME --keyring-backend $KEYRING_BACKEND)

out=$($BINARY tx wasm store ${HACK_ATOM} --from $KEY --output json --gas auto --gas-adjustment 2.3 --fees 100000000uluna --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not store contract" >&2
    echo $out >&2
    exit $code
fi
sleep 10
txhash=$(echo $out | jq -r '.txhash')
echo "$txhash"
code_id=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[1].attributes[1].value')
echo "code_id $code_id"

# instantiate stargate-tester
echo "... instantiates contract"
msg=$(jq -n '
{
   "verifier":"'$addr'",
   "beneficiary":"'$addr'"
}')
echo $msg
out=$($BINARY tx wasm instantiate $code_id "$msg" --from $KEY --output json --gas auto --gas-adjustment 2.3 --label "hackatom" --fees 20000000uluna --no-admin --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y)
code=$(echo $out | jq -r '.code')
if [ "$code" != "0" ]; then
    echo "... Could not instantiate contract" >&2
    echo $out >&2
    exit $code
fi
sleep 10

txhash=$(echo $out | jq -r '.txhash')
echo "$txhash"
contract_addr=$($BINARY q tx $txhash -o json | jq -r '.raw_log' | jq -r '.[0].events[0].attributes[0].value')
echo "contract_addr $contract_addr"

# query verifier
echo "... query verifier address"
msg='{"verifier":{}}'
out=$($BINARY query wasm contract-state smart $contract_addr $msg --output json)
echo $out

# execute release
echo "... execute release"
msg=$(jq -n '
{
   "nothing":{}
}')
echo $msg
out=$($BINARY tx wasm execute $contract_addr "$msg" --from $KEY --amount 1000000uluna --output json --gas auto --fees 400000000uluna --gas-adjustment 2.3 --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y -o json)
echo $out
txhash=$(echo $out | jq -r '.txhash')
echo $txhash
sleep 5
gasUsed=$($BINARY q tx --type=hash $txhash -o json | jq -r '.gas_used')
echo gas used $gasUsed before adding to tax exemption list

sleep 5
echo "... query tax-proceeds before adding to tax exemption list"
tax_proceeds_0=$($BINARY q treasury tax-proceeds)

# add contract_addr to burn tax exemption list
echo "add contract $contract_addr to burn tax exemption list"
out=$($BINARY tx gov submit-proposal add-burn-tax-exemption-address "$contract_addr" --title "burn tax exemption address" --description "burn tax exemption address"  --from $KEY --keyring-backend $KEYRING_BACKEND --chain-id $CHAIN_ID --home $HOME -y)
echo $out

sleep 5
out=$($BINARY tx gov deposit 1 "20000000${DENOM}" --from $KEY --keyring-backend $KEYRING_BACKEND --chain-id $CHAIN_ID --home $HOME -y)
echo $out

sleep 5
out=$($BINARY tx gov vote 1 yes --from $KEY --keyring-backend $KEYRING_BACKEND --chain-id $CHAIN_ID --home $HOME -y)
echo $out

sleep 5
out=$($BINARY tx gov vote 1 yes --from $KEY1 --keyring-backend $KEYRING_BACKEND --chain-id $CHAIN_ID --home $HOME -y)
echo $out

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

echo ""
echo "CHECK ADDRESS AFTER ADDING BURN TAX EXEMPTION LIST"
echo ""

# check burn tax exemption address
./_build/new/terrad q treasury burn-tax-exemption-list -o json | jq ".addresses"

# execute release
msg=$(jq -n '
{
   "nothing":{}
}')
echo $msg
out=$($BINARY tx wasm execute $contract_addr "$msg" --from $KEY --output json --gas auto --gas-adjustment 2.3 --chain-id $CHAIN_ID --home $HOME --keyring-backend $KEYRING_BACKEND -y -o json)
echo $out
txhash=$(echo $out | jq -r '.txhash')
echo $txhash
sleep 5
gasUsed=$($BINARY q tx --type=hash $txhash -o json | jq -r '.gas_used')
echo gas used $gasUsed after adding to tax exemption list

echo "... print tax-proceeds before adding to tax exemption list"
echo $tax_proceeds_0
echo "... query tax-proceeds after adding to tax exemption list, should be the same as before"
$BINARY q treasury tax-proceeds