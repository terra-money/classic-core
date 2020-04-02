package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Defines the prefix of each query path
const (
	QueryParameters = "parameters"
	QueryAuctions   = "auctions"
	QueryBids       = "bids"
	QueryRegistry   = "registry"
	QueryResolve    = "resolve"
	QueryReverse    = "reverse"
)

// QueryAuctionsParams defines the params for the following queries:
// - 'custom/nameservice/auctions'
type QueryAuctionsParams struct {
	NameHash NameHash
	Status   AuctionStatus
}

// QueryBidsParams defines the params for the following queries:
// - 'custom/nameservice/bids'
type QueryBidsParams struct {
	NameHash NameHash
	Bidder   sdk.AccAddress
}

// QueryRegistryParams defines the params for the following queries:
// - 'custom/nameservice/registry'
type QueryRegistryParams struct {
	NameHash NameHash
}

// QueryResolveParams defines the params for the following queries:
// - 'custom/nameservice/resolve'
type QueryResolveParams struct {
	NameHash      NameHash
	ChildNameHash NameHash
}

// QueryReverseParams defines the params for the following queries:
// - 'custom/nameservice/reverse'
type QueryReverseParams struct {
	Address sdk.AccAddress
}
