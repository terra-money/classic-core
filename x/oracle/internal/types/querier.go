package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines the prefix of each query path
const (
	QueryParameters       = "parameters"
	QueryPrice            = "price"
	QueryActives          = "actives"
	QueryPrevotes         = "prevotes"
	QueryVotes            = "votes"
	QueryFeederDelegation = "feederDelegation"
)

// QueryPriceParams defines the params for the following queries:
// - 'custom/oracle/price'
type QueryPriceParams struct {
	Denom string
}

// NewQueryPriceParams returns params for price query
func NewQueryPriceParams(denom string) QueryPriceParams {
	return QueryPriceParams{denom}
}

// QueryPrevotesParams defines the params for the following queries:
// - 'custom/oracle/prevotes'
type QueryPrevotesParams struct {
	Voter sdk.ValAddress
	Denom string
}

// NewQueryPrevotesParams returns params for price prevotes query
func NewQueryPrevotesParams(voter sdk.ValAddress, denom string) QueryPrevotesParams {
	return QueryPrevotesParams{voter, denom}
}

// QueryVotesParams defines the params for the following queries:
// - 'custom/oracle/votes'
type QueryVotesParams struct {
	Voter sdk.ValAddress
	Denom string
}

// NewQueryVotesParams returns params for price votes query
func NewQueryVotesParams(voter sdk.ValAddress, denom string) QueryVotesParams {
	return QueryVotesParams{voter, denom}
}

// QueryFeederDelegationParams defeins the params for the following queries:
// - 'custom/oracle/feederDelegation'
type QueryFeederDelegationParams struct {
	Validator sdk.ValAddress
}

// NewQueryFeederDelegationParams returns params for feeder delegation query
func NewQueryFeederDelegationParams(validator sdk.ValAddress) QueryFeederDelegationParams {
	return QueryFeederDelegationParams{validator}
}
