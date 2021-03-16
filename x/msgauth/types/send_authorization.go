package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// NewSendAuthorization return new SendAuthorization instance
func NewSendAuthorization(spendLimit sdk.Coins) *SendAuthorization {
	return &SendAuthorization{SpendLimit: spendLimit}
}

// MsgType return msg type of the authorization
func (authorization SendAuthorization) MsgType() string {
	return banktypes.TypeMsgSend
}

// Accept return whether the operation is allowed or not, and also
// returns the updated authorization.
func (authorization SendAuthorization) Accept(msg sdk.Msg, block tmproto.Header) (allow bool, updated AuthorizationI, delete bool) {
	switch msg := msg.(type) {
	case *banktypes.MsgSend:
		limitLeft, isNegative := authorization.SpendLimit.SafeSub(msg.Amount)
		if isNegative {
			return false, nil, false
		}
		if limitLeft.IsZero() {
			return true, nil, true
		}

		return true, &SendAuthorization{SpendLimit: limitLeft}, false
	}
	return false, nil, false
}
