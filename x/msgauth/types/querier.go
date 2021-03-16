package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines the prefix of each query path
const (
	QueryGrants    = "grants"
	QueryAllGrants = "all_grants"
)

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

// QueryAllGrantsParams defines the params for the following queries:
// - 'custom/msgauth/all_grants'
type QueryAllGrantsParams struct {
	Granter sdk.AccAddress `json:"granter"`
}

// NewQueryAllGrantsParams returns params for grant query
func NewQueryAllGrantsParams(granter sdk.AccAddress) QueryAllGrantsParams {
	return QueryAllGrantsParams{Granter: granter}
}
