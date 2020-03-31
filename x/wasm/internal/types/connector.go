package types

import (
	"fmt"
	"github.com/terra-project/core/x/market"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	wasmTypes "github.com/confio/go-cosmwasm/types"
)

// ParseMsgSend converts wasm msg to sdk.Msg
func ParseMsgSend(wasmMsg *wasmTypes.SendMsg) (msgSend bank.MsgSend, err sdk.Error) {
	fromAddr, stderr := sdk.AccAddressFromBech32(wasmMsg.FromAddress)
	if stderr != nil {
		err = sdk.ErrInvalidAddress(wasmMsg.FromAddress)
		return
	}
	toAddr, stderr := sdk.AccAddressFromBech32(wasmMsg.ToAddress)
	if stderr != nil {
		err = sdk.ErrInvalidAddress(wasmMsg.ToAddress)
		return
	}

	coins, err := ParseToCoins(wasmMsg.Amount)
	if err != nil {
		return
	}

	msgSend = bank.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      coins,
	}

	return
}

// ParseMsgSend converts wasm msg to sdk.Msg
func ParseMsgSwap(wasmMsg *wasmTypes.SwapMsg) (msgSend market.MsgSwap, err sdk.Error) {
	traderAddr, stderr := sdk.AccAddressFromBech32(wasmMsg.TraderAddress)
	if stderr != nil {
		err = sdk.ErrInvalidAddress(wasmMsg.TraderAddress)
		return
	}

	offerCoin, err := ParseToCoin(wasmMsg.OfferCoin)
	if err != nil {
		return
	}

	msgSend = market.MsgSwap{
		Trader:    traderAddr,
		OfferCoin: offerCoin,
		AskDenom:  wasmMsg.AskDenom,
	}

	return
}

// ParseOpaqueMsg decodes msg.Data to an sdk.Msg using amino json encoding.
func ParseOpaqueMsg(cdc *codec.Codec, wasmMsg *wasmTypes.OpaqueMsg) (sdk.Msg, sdk.Error) {
	// until more is changes, format is amino json encoding, wrapped base64
	var msg sdk.Msg
	err := cdc.UnmarshalJSON(wasmMsg.Data, &msg)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("failed to parse opaque msg to sdk msg", err.Error()))
	}

	return msg, nil
}

// ToWasmMsg encodes an sdk.Msg using amino json encoding.
// Then wraps it as an opaque message
func ToWasmMsg(cdc *codec.Codec, msg sdk.Msg) (wasmTypes.CosmosMsg, sdk.Error) {
	opaqueBz, err := cdc.MarshalJSON(msg)
	if err != nil {
		return wasmTypes.CosmosMsg{}, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	wasmMsg := wasmTypes.CosmosMsg{
		Opaque: &wasmTypes.OpaqueMsg{
			Data: opaqueBz,
		},
	}

	return wasmMsg, nil
}

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
