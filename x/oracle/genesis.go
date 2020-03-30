package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for delegatorBechAddr, delegate := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegate)
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

	for delegatorBechAddr, delegate := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegate)
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

	for _, aggregatePrevote := range data.AggregateExchangeRatePrevotes {
		keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)
	}

	for _, aggregateVote := range data.AggregateExchangeRateVotes {
		keeper.AddAggregateExchangeRateVote(ctx, aggregateVote)
	}

	if len(data.TobinTaxes) > 0 {
		for denom, tobinTax := range data.TobinTaxes {
			keeper.SetTobinTax(ctx, denom, tobinTax)
		}
	} else {
		for _, item := range data.Params.Whitelist {
			keeper.SetTobinTax(ctx, item.Name, item.TobinTax)
		}
	}

	keeper.SetParams(ctx, data.Params)

	// check if the module account exists
	moduleAcc := keeper.GetOracleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	feederDelegations := make(map[string]sdk.AccAddress)
	keeper.IterateOracleDelegates(ctx, func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool) {
		bechAddr := delegator.String()
		feederDelegations[bechAddr] = delegate
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

	var aggregateExchangeRatePrevotes []AggregateExchangeRatePrevote
	keeper.IterateAggregateExchangeRatePrevotes(ctx, func(aggregatePrevote AggregateExchangeRatePrevote) (stop bool) {
		aggregateExchangeRatePrevotes = append(aggregateExchangeRatePrevotes, aggregatePrevote)
		return false
	})

	var aggregateExchangeRateVotes []AggregateExchangeRateVote
	keeper.IterateAggregateExchangeRateVotes(ctx, func(aggregateVote AggregateExchangeRateVote) bool {
		aggregateExchangeRateVotes = append(aggregateExchangeRateVotes, aggregateVote)
		return false
	})

	tobinTaxes := make(map[string]sdk.Dec)
	keeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		tobinTaxes[denom] = tobinTax
		return false
	})

	return NewGenesisState(params, exchangeRatePrevotes, exchangeRateVotes, rates, feederDelegations, missCounters, aggregateExchangeRatePrevotes, aggregateExchangeRateVotes, tobinTaxes)
}
