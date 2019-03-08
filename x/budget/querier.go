package budget

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryPrice  = "price"
	QueryVotes  = "votes"
	QueryActive = "active"
	QueryParams = "params"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPrice:
			return queryPrice(ctx, path[1:], req, keeper)
		case QueryActive:
			return queryActive(ctx, req, keeper)
		case QueryVotes:
			return queryVotes(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

// nolint: unparam
func queryPrice(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]

	price, err := keeper.GetPrice(ctx, denom)

	if err != nil {
		return []byte{}, ErrUnknownDenomination(DefaultCodespace, denom)
	}

	return price.Bytes(), nil
}

// nolint: unparam
func queryActive(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denoms := []string{}

	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixPrice)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var denom string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Key(), &denom)
		denoms = append(denoms, denom)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, denoms)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

// Params for query 'custom/oracle/votes'
type QueryVoteParams struct {
	Voter sdk.AccAddress
	Denom string
}

// creates a new instance of QueryVoteParams
func NewQueryVoteParams(voter sdk.AccAddress, denom string) QueryVoteParams {
	return QueryVoteParams{
		Voter: voter,
		Denom: denom,
	}
}

// nolint: unparam
func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredVotes := PriceBallot{}
	votes := keeper.getVotes(ctx)

	for _, ballot := range votes {
		for _, vote := range ballot {
			if len(params.Denom) != 0 && len(params.Voter) != 0 {
				if vote.Denom == params.Denom && vote.Voter.Equals(params.Voter) {
					filteredVotes = append(filteredVotes, vote)
				}

			} else if len(params.Denom) != 0 {
				if vote.Denom == params.Denom {
					filteredVotes = append(filteredVotes, vote)
				}
			} else if len(params.Voter) != 0 {
				if vote.Voter.Equals(params.Voter) {
					filteredVotes = append(filteredVotes, vote)
				}
			} else {
				myVote := keeper.getVote(ctx, params.Denom, params.Voter)
				filteredVotes = append(filteredVotes, myVote)
			}
		}

	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, filteredVotes)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
