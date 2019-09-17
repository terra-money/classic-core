package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryParameters       = "parameters"
	QueryPrice            = "price"
	QueryActives          = "actives"
	QueryPrevotes         = "prevotes"
	QueryVotes            = "votes"
	QueryFeederDelegation = "feederDelegation"
	QueryVotingInfo       = "signingInfo"
	QueryVotingInfos      = "signingInfos"
)

// QueryPriceParams defines the params for the following queries:
// - 'custom/oracle/price'
type QueryPriceParams struct {
	Denom string
}

func NewQueryPriceParams(denom string) QueryPriceParams {
	return QueryPriceParams{denom}
}

// QueryPrevotesParams defines the params for the following queries:
// - 'custom/oracle/prevotes'
type QueryPrevotesParams struct {
	Voter sdk.ValAddress
	Denom string
}

func NewQueryPrevotesParams(voter sdk.ValAddress, denom string) QueryPrevotesParams {
	return QueryPrevotesParams{voter, denom}
}

// QueryVotesParams defines the params for the following queries:
// - 'custom/oracle/votes'
type QueryVotesParams struct {
	Voter sdk.ValAddress
	Denom string
}

func NewQueryVotesParams(voter sdk.ValAddress, denom string) QueryVotesParams {
	return QueryVotesParams{voter, denom}
}

// QueryFeederDelegationParams defeins the params for the following queries:
// - 'custom/oracle/feederDelegation'
type QueryFeederDelegationParams struct {
	Validator sdk.ValAddress
}

func NewQueryFeederDelegationParams(validator sdk.ValAddress) QueryFeederDelegationParams {
	return QueryFeederDelegationParams{validator}
}

// QueryVotingInfoParams defines the params for the following queries:
// - 'custom/oracle/votingInfo'
type QueryVotingInfoParams struct {
	ValAddress sdk.ValAddress
}

func NewQueryVotingInfoParams(valAddr sdk.ValAddress) QueryVotingInfoParams {
	return QueryVotingInfoParams{valAddr}
}

// QueryVotingInfosParams defines the params for the following queries:
// - 'custom/oracle/votingInfos'
type QueryVotingInfosParams struct {
	Page, Limit int
}

func NewQueryVotingInfosParams(page, limit int) QueryVotingInfosParams {
	return QueryVotingInfosParams{page, limit}
}
