package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenericAuthorization grants the permission to execute any transaction of the provided
// msg type without restrictions
type GenericAuthorization struct {
	// GrantMsgType is the type of Msg this capability grant allows
	GrantMsgType string `json:"grant_msg_type"`
}

// NewGenericAuthorization returns new GenericAuthorization instantce
func NewGenericAuthorization(msgType string) GenericAuthorization {
	return GenericAuthorization{GrantMsgType: msgType}
}

// MsgType implement Authorization
func (ga GenericAuthorization) MsgType() string {
	return ga.GrantMsgType
}

// Accept implement Authorization
func (ga GenericAuthorization) Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Authorization, delete bool) {
	return true, ga, false
}
