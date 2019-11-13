package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// InitGenesis initializes default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	keeper.SetTaxRate(ctx, data.TaxRate)
	keeper.SetRewardWeight(ctx, data.RewardWeight)
	keeper.SetEpochInitialIssuance(ctx, data.EpochInitialIssuance)
	keeper.SetTaxProceeds(ctx, data.TaxProceed)

	// store tax caps
	for denom, taxCap := range data.TaxCaps {
		keeper.SetTaxCap(ctx, denom, taxCap)
	}

	fmt.Println(data.MRs)

	for epoch, MR := range data.MRs {
		keeper.SetMR(ctx, int64(epoch), MR)
	}
	for epoch, SR := range data.SRs {
		keeper.SetSR(ctx, int64(epoch), SR)
	}
	for epoch, TRL := range data.TRLs {
		keeper.SetTRL(ctx, int64(epoch), TRL)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)

	taxRate := keeper.GetTaxRate(ctx)
	rewardWeight := keeper.GetRewardWeight(ctx)
	taxProceeds := keeper.PeekTaxProceeds(ctx)
	epochInitialIssuance := keeper.GetEpochInitialIssuance(ctx)

	taxCaps := make(map[string]sdk.Int)
	keeper.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps[denom] = taxCap
		return false
	})

	var MRs []sdk.Dec
	var SRs []sdk.Dec
	var TRLs []sdk.Dec

	curEpoch := core.GetEpoch(ctx)
	for e := int64(0); e < curEpoch ||
		(e == curEpoch && core.IsPeriodLastBlock(ctx, core.BlocksPerEpoch)); e++ {

		MRs = append(MRs, keeper.GetMR(ctx, e))
		SRs = append(SRs, keeper.GetSR(ctx, e))
		TRLs = append(TRLs, keeper.GetTRL(ctx, e))
	}

	return NewGenesisState(params, taxRate, rewardWeight,
		taxCaps, taxProceeds, epochInitialIssuance, MRs, SRs, TRLs)
}
