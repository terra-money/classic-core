package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines the prefix of each query path
const (
	QueryGrant  = "grant"
	QueryGrants = "grants"
)

// QueryGrantParams defines the params for the following queries:
// - 'custom/msgauth/grant'
type QueryGrantParams struct {
	Granter sdk.AccAddress `json:"granter"`
	Grantee sdk.AccAddress `json:"grantee"`
	MsgType string         `json:"msg_type"`
}

// NewQueryGrantParams returns params for grant query
func NewQueryGrantParams(granter sdk.AccAddress, grantee sdk.AccAddress, msgType string) QueryGrantParams {
	return QueryGrantParams{Granter: granter, Grantee: grantee, MsgType: msgType}
}

// QueryGrantsParams defines the params for the following queries:
// - 'custom/msgauth/grants'
type QueryGrantsParams struct {
	Granter sdk.AccAddress `json:"granter"`
	Grantee sdk.AccAddress `json:"grantee"`
}

// NewQueryGrantsParams returns params for grant query
func NewQueryGrantsParams(granter sdk.AccAddress, grantee sdk.AccAddress) QueryGrantsParams {
	return QueryGrantsParams{Granter: granter, Grantee: grantee}
}
