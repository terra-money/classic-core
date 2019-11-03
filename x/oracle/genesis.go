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

	for _, prevote := range data.ExchangeRatePrevotes {
		keeper.AddExchangeRatePrevote(ctx, prevote)
	}

	for _, vote := range data.ExchangeRateVotes {
		keeper.AddExchangeRateVote(ctx, vote)
	}

	for denom, rate := range data.ExchangeRates {
		keeper.SetLunaExchangeRate(ctx, denom, rate)
	}

	for delegatorBechAddr, delegatee := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegatee)
	}

	for _, prevote := range data.ExchangeRatePrevotes {
		keeper.AddExchangeRatePrevote(ctx, prevote)
	}

	for _, vote := range data.ExchangeRateVotes {
		keeper.AddExchangeRateVote(ctx, vote)
	}

	for denom, rate := range data.ExchangeRates {
		keeper.SetLunaExchangeRate(ctx, denom, rate)
	}

	for operatorBechAddr, missCounter := range data.MissCounters {
		operator, err := sdk.ValAddressFromBech32(operatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetMissCounter(ctx, operator, missCounter)
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

	var exchangeRatePrevotes []ExchangeRatePrevote
	keeper.IterateExchangeRatePrevotes(ctx, func(prevote ExchangeRatePrevote) (stop bool) {
		exchangeRatePrevotes = append(exchangeRatePrevotes, prevote)
		return false
	})

	var exchangeRateVotes []ExchangeRateVote
	keeper.IterateExchangeRateVotes(ctx, func(vote ExchangeRateVote) (stop bool) {
		exchangeRateVotes = append(exchangeRateVotes, vote)
		return false
	})

	rates := make(map[string]sdk.Dec)
	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates[denom] = rate
		return false
	})

	missCounters := make(map[string]int64)
	keeper.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter int64) (stop bool) {
		missCounters[operator.String()] = missCounter
		return false
	})

	return NewGenesisState(params, exchangeRatePrevotes, exchangeRateVotes, rates, feederDelegations, missCounters)
}
