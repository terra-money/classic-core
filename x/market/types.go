package market

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//--------------------------------------------------------
//--------------------------------------------------------

// ReserveParams defines the basic properties of a staking proposal
type ReserveParams struct {
	Target Int
	Cap    Int
}

// const (
// 	minReserveRatio = 1.2
// 	maxReserveRatio = 1.5
// )

// func (r *ReserveParams) isOverCapitalized(terraCap Int, terraPrice sdk.Rat, lunaPrice sdk.Rat) {
// 	return terraCap*terraPrice*maxReserveRatio < r.Current*lunaPrice
// }

// func (r *ReserveParams) isUnderCapitalized(terraCap Int, terraPrice sdk.Rat, lunaPrice sdk.Rat) {
// 	return terraCap*terraPrice*minReserveRatio > r.Current*lunaPrice
// }

//--------------------------------------------------------
//--------------------------------------------------------

// SwapMsg defines the msg of a trader containing terra coin to be
// swapped with luna coin, or luna coin to be swapped with the terra coin
type SwapMsg struct {
	Trader sdk.AccAddress // Address of the trader
	Coin   sdk.Coin       // Coin to be swapped
}

// NewVoteMsg creates a VoteMsg instance
func NewSwapMsg(traderAddress sdk.AccAddress, coin sdk.Coin) SwapMsg {
	// coin must not be nil, and only Terra and Luna can be swapped
	if coin == nil || !isValidCoin(coin) {
		return nil
	}

	return SwapMsg{
		Trader: traderAddress,
		Coin:   coin,
	}
}

// Type Implements Msg
func (msg SwapMsg) Type() string {
	return "market"
}

// Get Implements Msg
func (msg SwapMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes Implements Msg
func (msg SwapMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners Implements Msg
func (msg SwapMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

func isValidCoin(coin sdk.Coin) bool {
	denoms := []string{"terra", "luna"}
	for _, value := range options {
		if value == coin.Denom {
			return true
		}
	}
	return false
}

// ValidateBasic Implements Msg
func (msg SwapMsg) ValidateBasic() sdk.Error {
	if len(msg.Trader) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Trader.String())
	}
	if !isValidCoin(msg.Coin) {
		return ErrInvalidOption("Invalid coin: " + msg.Coin)
	}

	return nil
}

// String Implements Msg
func (msg SwapMsg) String() string {
	return fmt.Sprintf("SwapMsg{%v, %v}", msg.Trader, msg.Coin)
}
