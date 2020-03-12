package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryAuctions:
			return queryAuctions(ctx, req, keeper)
		case types.QueryBids:
			return queryBids(ctx, req, keeper)
		case types.QueryRegistry:
			return queryRegistry(ctx, req, keeper)
		case types.QueryResolve:
			return queryResolve(ctx, req, keeper)
		case types.QueryReverse:
			return queryReverse(ctx, req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryAuctions(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryAuctionsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	var auctions types.Auctions
	if len(params.NameHash) != 0 {
		auction, err := keeper.GetAuction(ctx, params.NameHash)
		if err == nil {
			auctions = append(auctions, auction)
		}
	} else {
		keeper.IterateAuction(ctx, func(nameHash types.NameHash, auction types.Auction) bool {
			if params.Status == types.AuctionStatusNil || auction.Status == params.Status {
				auctions = append(auctions, auction)
			}
			return false
		})
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, auctions)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryBids(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryBidsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if params.NameHash.Empty() {
		return nil, sdk.ErrUnknownRequest("invalid params; need name hash")
	}

	var bids []types.Bid
	if !params.Bidder.Empty() {
		bid, err := keeper.GetBid(ctx, params.NameHash, params.Bidder)
		if err == nil {
			bids = append(bids, bid)
		}
	} else {
		keeper.IterateBid(ctx, params.NameHash, func(_ types.NameHash, bid types.Bid) bool {
			bids = append(bids, bid)
			return false
		})
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, bids)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryRegistry(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryRegistryParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if params.NameHash.Empty() {
		return nil, sdk.ErrUnknownRequest("invalid params; need name hash")
	}

	registry, err2 := keeper.GetRegistry(ctx, params.NameHash)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, registry)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryResolve(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryResolveParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if params.NameHash.Empty() {
		return nil, sdk.ErrUnknownRequest("invalid params; need both parent and child name hash")
	}

	resolve, err2 := keeper.GetResolve(ctx, params.NameHash, params.ChildNameHash)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, resolve)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryReverse(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryReverseParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if params.Address.Empty() {
		return nil, sdk.ErrUnknownRequest("invalid params, need address")
	}

	nameHash, err2 := keeper.GetReverseResolve(ctx, params.Address)
	if err2 != nil {
		return nil, err2
	}

	registry, err2 := keeper.GetRegistry(ctx, nameHash)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, registry)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
