package types

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenericAuthorization returns new GenericAuthorization instance
func NewGenericAuthorization(msgType string) *GenericAuthorization {
	return &GenericAuthorization{GrantMsgType: msgType}
}

// MsgType implement Authorization
func (ga GenericAuthorization) MsgType() string {
	return ga.GrantMsgType
}

// Accept implement Authorization
func (ga GenericAuthorization) Accept(msg sdk.Msg, blocktime tmproto.Header) (allow bool, updated AuthorizationI, delete bool) {
	return true, &ga, false
}
