package types

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

// DefaultFeatures - Cosmwasm feature
const DefaultFeatures = "stargate,staking,terra"

// ParseEvents converts wasm EventAttributes and Events into an sdk.Events
func ParseEvents(
	contractAddr sdk.AccAddress,
	attributes wasmvmtypes.EventAttributes,
	events wasmvmtypes.Events,
) (sdk.Events, error) {
	if len(attributes) == 0 && len(events) == 0 {
		return nil, nil
	}

	var sdkEvents sdk.Events

	if len(attributes) != 0 {
		sdkEvent, err := buildEvent(EventTypeWasmPrefix, contractAddr, attributes)
		if err != nil {
			return nil, err
		}

		sdkEvents = sdkEvents.AppendEvent(*sdkEvent)

		// Deprecated: from_contract
		sdkEvent.Type = EventTypeFromContract
		sdkEvents = sdkEvents.AppendEvent(*sdkEvent)
	}

	// append wasm prefix for the events
	for _, event := range events {
		sdkEvent, err := buildEvent(fmt.Sprintf("%s-%s", EventTypeWasmPrefix, event.Type), contractAddr, event.Attributes)
		if err != nil {
			return nil, err
		}

		sdkEvents = sdkEvents.AppendEvent(*sdkEvent)
	}

	return sdkEvents, nil
}

func buildEvent(
	eventType string,
	contractAddr sdk.AccAddress,
	attributes wasmvmtypes.EventAttributes,
) (*sdk.Event, error) {
	if len(attributes) == 0 {
		return nil, nil
	}

	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddress, contractAddr.String())}
	for _, l := range attributes {
		// and reserve the contract_address key for our use (not contract)
		if l.Key != AttributeKeyContractAddress {
			attr := sdk.NewAttribute(l.Key, l.Value)
			attrs = append(attrs, attr)
		}
	}

	event := sdk.NewEvent(eventType, attrs...)
	return &event, nil
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
// Deprecated `from_contract` will be excluded from the events
func EncodeSdkEvents(events []sdk.Event) []wasmvmtypes.Event {
	res := make([]wasmvmtypes.Event, len(events))
	for i, ev := range events {
		// Deprecated: from_contract
		if ev.Type == EventTypeFromContract {
			continue
		}

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
