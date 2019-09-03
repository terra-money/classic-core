package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetTaxRate(ctx, data.TaxRate)
	keeper.SetRewardWeight(ctx, data.RewardWeight)

	// store tax cap for SDT & LUNA(no tax)
	keeper.SetTaxCap(ctx, data.Params.TaxPolicy.Cap.Denom, data.Params.TaxPolicy.Cap.Amount)
	keeper.SetTaxCap(ctx, core.MicroLunaDenom, sdk.ZeroInt())
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	taxRate := keeper.GetTaxRate(ctx, core.GetEpoch(ctx))
	rewardWeight := keeper.GetRewardWeight(ctx, core.GetEpoch(ctx))
	return NewGenesisState(params, taxRate, rewardWeight)
}
