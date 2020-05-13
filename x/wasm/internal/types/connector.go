package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
)

// ParseResult converts wasm result to sdk.Result
func ParseResult(wasmResult *wasmTypes.Result, contractAddr sdk.AccAddress) sdk.Result {
	var events []sdk.Event
	if len(wasmResult.Log) > 0 {
		// we always tag with the contract address issuing this event
		attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddress, contractAddr.String())}
		for _, l := range wasmResult.Log {
			// and reserve the contract_address key for our use (not contract)
			if l.Key != AttributeKeyContractAddress {
				attr := sdk.NewAttribute(l.Key, l.Value)
				attrs = append(attrs, attr)
			}
		}

		events = []sdk.Event{sdk.NewEvent(EventTypeExecuteContract, attrs...)}
	}

	return sdk.Result{
		Data:   []byte(wasmResult.Data),
		Events: events,
	}
}

// ParseToCoin converts wasm coin to sdk.Coin
func ParseToCoin(wasmCoin wasmTypes.Coin) (coin sdk.Coin, err sdk.Error) {
	amount, ok := sdk.NewIntFromString(wasmCoin.Amount)
	if !ok {
		err = sdk.ErrInvalidCoins(fmt.Sprintf("Failed to parse %s", coin.Amount))
		return
	}

	coin = sdk.Coin{
		Denom:  wasmCoin.Denom,
		Amount: amount,
	}
	return
}

// ParseToCoins converts wasm coins to sdk.Coins
func ParseToCoins(wasmCoins []wasmTypes.Coin) (coins sdk.Coins, err sdk.Error) {
	for _, coin := range wasmCoins {
		c, err := ParseToCoin(coin)
		if err != nil {
			return nil, err
		}

		coins = append(coins, c)
	}
	return
}
