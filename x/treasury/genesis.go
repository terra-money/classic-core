package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// InitGenesis initializes default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for epoch, taxRate := range data.TaxRates {
		keeper.SetTaxRate(ctx, int64(epoch), taxRate)
	}

	for epoch, rewardWeight := range data.RewardWeights {
		keeper.SetRewardWeight(ctx, int64(epoch), rewardWeight)
	}

	// store tax caps
	for denom, taxCap := range data.TaxCaps {
		keeper.SetTaxCap(ctx, denom, taxCap)
	}

	for epoch, historicalIssuance := range data.HistoricalIssuances {
		keeper.SetHistoricalIssuance(ctx, int64(epoch), historicalIssuance)
	}

	for epoch, taxProceed := range data.TaxProceeds {
		keeper.SetTaxProceeds(ctx, int64(epoch), taxProceed)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)

	taxCaps := make(map[string]sdk.Int)
	keeper.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps[denom] = taxCap
		return false
	})

	var taxRates []sdk.Dec
	var rewardWeights []sdk.Dec
	var taxProceeds []sdk.Coins
	var historicalIssuancees []sdk.Coins
	for e := int64(0); e <= core.GetEpoch(ctx); e++ {
		taxRates = append(taxRates, keeper.GetTaxRate(ctx, e))
		rewardWeights = append(rewardWeights, keeper.GetRewardWeight(ctx, e))
		taxProceeds = append(taxProceeds, keeper.PeekTaxProceeds(ctx, e))
		historicalIssuancees = append(historicalIssuancees, keeper.GetHistoricalIssuance(ctx, e))
	}

	return NewGenesisState(params, taxRates, rewardWeights, taxCaps, taxProceeds, historicalIssuancees)
}
