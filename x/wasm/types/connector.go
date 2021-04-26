package types

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

// DefaultFeatures - Cosmwasm feature
const DefaultFeatures = "stargate,staking,terra"

// ValidateAndParseEvents converts wasm LogAttributes into an sdk.Events (with 0 or 1 elements)
func ValidateAndParseEvents(contractAddr sdk.AccAddress, params EventParams, attributes ...wasmvmtypes.EventAttribute) (sdk.Events, error) {
	if len(attributes) == 0 {
		return nil, nil
	}

	if len(attributes) > int(params.MaxAttributeNum) {
		return nil, ErrExceedMaxContractEventAttributeNum
	}

	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddress, contractAddr.String())}
	for _, l := range attributes {
		if len(l.Key) > int(params.MaxAttributeKeyLength) {
			return nil, ErrExceedMaxContractEventAttributeKeyLength
		}

		if len(l.Value) > int(params.MaxAttributeValueLength) {
			return nil, ErrExceedMaxContractEventAttributeValueLength
		}

		// and reserve the contract_address key for our use (not contract)
		if l.Key != AttributeKeyContractAddress {
			attr := sdk.NewAttribute(l.Key, l.Value)
			attrs = append(attrs, attr)
		}
	}

	return sdk.Events{sdk.NewEvent(EventTypeFromContract, attrs...)}, nil
}

// ParseToCoin converts wasm coin to sdk.Coin
func ParseToCoin(wasmCoin wasmvmtypes.Coin) (coin sdk.Coin, err error) {
	amount, ok := sdk.NewIntFromString(wasmCoin.Amount)
	if !ok {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("Failed to parse %s", coin.Amount))
		return
	}

	coin = sdk.Coin{
		Denom:  wasmCoin.Denom,
		Amount: amount,
	}
	return
}

// ParseToCoins converts wasm coins to sdk.Coins
func ParseToCoins(wasmCoins []wasmvmtypes.Coin) (coins sdk.Coins, err error) {
	for _, coin := range wasmCoins {
		c, err := ParseToCoin(coin)
		if err != nil {
			return nil, err
		}

		coins = append(coins, c)
	}
	return
}

// EncodeSdkCoin - encode sdk coin to wasm coin
func EncodeSdkCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}

// EncodeSdkCoins - encode sdk coins to wasm coins
func EncodeSdkCoins(coins sdk.Coins) wasmvmtypes.Coins {
	encodedCoins := make(wasmvmtypes.Coins, len(coins))
	for i, c := range coins {
		encodedCoins[i] = EncodeSdkCoin(c)
	}
	return encodedCoins
}

// EncodeSdkEvents - encode sdk events to wasm events
func EncodeSdkEvents(events []sdk.Event) []wasmvmtypes.Event {
	res := make([]wasmvmtypes.Event, len(events))
	for i, ev := range events {
		res[i] = wasmvmtypes.Event{
			Type:       ev.Type,
			Attributes: encodeSdkAttributes(ev.Attributes),
		}
	}
	return res
}

func encodeSdkAttributes(attrs []abci.EventAttribute) []wasmvmtypes.EventAttribute {
	res := make([]wasmvmtypes.EventAttribute, len(attrs))
	for i, attr := range attrs {
		res[i] = wasmvmtypes.EventAttribute{
			Key:   string(attr.Key),
			Value: string(attr.Value),
		}
	}
	return res
}

// ConvertWasmIBCTimeoutHeightToCosmosHeight convert wasm types to cosmos type
func ConvertWasmIBCTimeoutHeightToCosmosHeight(ibcTimeoutBlock *wasmvmtypes.IBCTimeoutBlock) ibcclienttypes.Height {
	if ibcTimeoutBlock == nil {
		return ibcclienttypes.NewHeight(0, 0)
	}
	return ibcclienttypes.NewHeight(ibcTimeoutBlock.Revision, ibcTimeoutBlock.Height)
}

// ConvertWasmIBCTimeoutTimestampToCosmosTimestamp convert wasm types to cosmos type
func ConvertWasmIBCTimeoutTimestampToCosmosTimestamp(timestamp *uint64) uint64 {
	if timestamp == nil {
		return 0
	}
	return *timestamp
}

// Messenger is an extension point for custom wasmVM message handling
type Messenger interface {
	// DispatchMessage encodes the wasmVM message and dispatches it.
	DispatchMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data []byte, err error)
}
