package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// InitGenesis initializes default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	keeper.SetTaxRate(ctx, data.TaxRate)
	keeper.SetRewardWeight(ctx, data.RewardWeight)
	keeper.SetEpochTaxProceeds(ctx, data.TaxProceed)

	// If EpochInitialIssuance is empty, we use current supply as epoch initial issuance
	if data.EpochInitialIssuance.IsZero() {
		keeper.RecordEpochInitialIssuance(ctx)
	} else {
		keeper.SetEpochInitialIssuance(ctx, data.EpochInitialIssuance)
	}

	// store tax caps
	for denom, taxCap := range data.TaxCaps {
		keeper.SetTaxCap(ctx, denom, taxCap)
	}

	// store cumulated block height of past chains
	keeper.SetCumulatedHeight(ctx, data.CumulatedHeight)

	for epoch, TR := range data.TRs {
		keeper.SetTR(ctx, int64(epoch), TR)
	}
	for epoch, SR := range data.SRs {
		keeper.SetSR(ctx, int64(epoch), SR)
	}
	for epoch, TSL := range data.TSLs {
		keeper.SetTSL(ctx, int64(epoch), TSL)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)

	taxRate := keeper.GetTaxRate(ctx)
	rewardWeight := keeper.GetRewardWeight(ctx)
	taxProceeds := keeper.PeekEpochTaxProceeds(ctx)
	epochInitialIssuance := keeper.GetEpochInitialIssuance(ctx)

	taxCaps := make(map[string]sdk.Int)
	keeper.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps[denom] = taxCap
		return false
	})

	cumulatedHeight := keeper.GetCumulatedHeight(ctx)

	var TRs []sdk.Dec
	var SRs []sdk.Dec
	var TSLs []sdk.Int

	curEpoch := keeper.GetEpoch(ctx)
	for e := int64(0); e < curEpoch ||
		(e == curEpoch && core.IsPeriodLastBlock(ctx, core.BlocksPerWeek)); e++ {

		TRs = append(TRs, keeper.GetTR(ctx, e))
		SRs = append(SRs, keeper.GetSR(ctx, e))
		TSLs = append(TSLs, keeper.GetTSL(ctx, e))
	}

	return NewGenesisState(params, taxRate, rewardWeight,
		taxCaps, taxProceeds, epochInitialIssuance,
		cumulatedHeight, TRs, SRs, TSLs)
}
