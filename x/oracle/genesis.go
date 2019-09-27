package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for addr, info := range data.VotingInfos {
		address, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		keeper.SetVotingInfo(ctx, address, info)
	}

	for addr, array := range data.MissedVotes {
		address, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		for _, missed := range array {
			keeper.SetMissedVoteBitArray(ctx, address, missed.Index, missed.Missed)
		}
	}

	for delegatorBechAddr, delegatee := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetFeedDelegate(ctx, delegator, delegatee)
	}

	for _, prevote := range data.PricePrevotes {
		keeper.AddPrevote(ctx, prevote)
	}

	for _, vote := range data.PriceVotes {
		keeper.AddVote(ctx, vote)
	}

	for denom, price := range data.Prices {
		keeper.SetLunaPrice(ctx, denom, price)
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	votingInfos := make(map[string]VotingInfo)
	missedVotes := make(map[string][]MissedVote)
	keeper.IterateVotingInfos(ctx, func(info VotingInfo) (stop bool) {
		bechAddr := info.Address.String()

		votingInfos[bechAddr] = info
		localMissedVotes := []MissedVote{}

		keeper.IterateMissedVoteBitArray(ctx, info.Address, func(index int64, missed bool) (stop bool) {
			localMissedVotes = append(localMissedVotes, NewMissedVote(index, missed))
			return false
		})
		missedVotes[bechAddr] = localMissedVotes

		return false
	})

	feederDelegations := make(map[string]sdk.AccAddress)
	keeper.IterateFeederDelegations(ctx, func(delegator sdk.ValAddress, delegatee sdk.AccAddress) (stop bool) {
		bechAddr := delegator.String()
		feederDelegations[bechAddr] = delegatee
		return false
	})

	var pricePrevotes []PricePrevote
	keeper.IteratePrevotes(ctx, func(prevote PricePrevote) (stop bool) {
		pricePrevotes = append(pricePrevotes, prevote)
		return false
	})

	var priceVotes []PriceVote
	keeper.IterateVotes(ctx, func(vote PriceVote) (stop bool) {
		priceVotes = append(priceVotes, vote)
		return false
	})

	prices := make(map[string]sdk.Dec)
	keeper.IterateLunaPrices(ctx, func(denom string, price sdk.Dec) bool {
		prices[denom] = price
		return false
	})

	return NewGenesisState(params, pricePrevotes, priceVotes, prices, votingInfos, missedVotes, feederDelegations)
}
