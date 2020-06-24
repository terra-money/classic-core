package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
)

type SendAuthorization struct {
	// SpendLimit specifies the maximum amount of tokens that can be spent
	// by this authorization and will be updated as tokens are spent. If it is
	// empty, there is no spend limit and any amount of coins can be spent.
	SpendLimit sdk.Coins `json:"spend_limit"`
}

// NewSendAuthorization return new SendAhtorization instance
func NewSendAuthorization(spendLimit sdk.Coins) SendAuthorization {
	return SendAuthorization{SpendLimit: spendLimit}
}

func (authorization SendAuthorization) MsgType() string {
	return bank.MsgSend{}.Type()
}

func (authorization SendAuthorization) Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Authorization, delete bool) {
	switch msg := msg.(type) {
	case bank.MsgSend:
		limitLeft, isNegative := authorization.SpendLimit.SafeSub(msg.Amount)
		if isNegative {
			return false, nil, false
		}
		if limitLeft.IsZero() {
			return true, nil, true
		}

		return true, SendAuthorization{SpendLimit: limitLeft}, false
	}
	return false, nil, false
}
