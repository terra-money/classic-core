package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for delegatorBechAddr, delegatee := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegatee)
	}

	for _, prevote := range data.Prevotes {
		keeper.AddPrevote(ctx, prevote)
	}

	for _, vote := range data.Votes {
		keeper.AddVote(ctx, vote)
	}

	for denom, exchangeRate := range data.ExchangeRates {
		keeper.SetLunaExchangeRate(ctx, denom, exchangeRate)
	}

	for delegatorBechAddr, delegatee := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegatee)
	}

	for _, prevote := range data.Prevotes {
		keeper.AddPrevote(ctx, prevote)
	}

	for _, vote := range data.Votes {
		keeper.AddVote(ctx, vote)
	}

	for denom, exchangeRate := range data.ExchangeRates {
		keeper.SetLunaExchangeRate(ctx, denom, exchangeRate)
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	feederDelegations := make(map[string]sdk.AccAddress)
	keeper.IterateOracleDelegates(ctx, func(delegator sdk.ValAddress, delegatee sdk.AccAddress) (stop bool) {
		bechAddr := delegator.String()
		feederDelegations[bechAddr] = delegatee
		return false
	})

	var Prevotes []Prevote
	keeper.IteratePrevotes(ctx, func(prevote Prevote) (stop bool) {
		Prevotes = append(Prevotes, prevote)
		return false
	})

	var Votes []Vote
	keeper.IterateVotes(ctx, func(vote Vote) (stop bool) {
		Votes = append(Votes, vote)
		return false
	})

	exchangeRates := make(map[string]sdk.Dec)
	keeper.IterateLunaExchangeRates(ctx, func(denom string, exchangeRate sdk.Dec) bool {
		exchangeRates[denom] = exchangeRate
		return false
	})

	return NewGenesisState(params, Prevotes, Votes, exchangeRates, feederDelegations)
}
