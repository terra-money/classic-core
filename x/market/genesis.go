package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetBasePool(ctx, data.BasePool)
	keeper.SetTerraPool(ctx, data.TerraPool)
	keeper.SetLastUpdateHeight(ctx, data.LastUpdateHeight)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	basePool := keeper.GetBasePool(ctx)
	terraPool := keeper.GetTerraPool(ctx)
	lastUpdateHeight := keeper.GetLastUpdateHeight(ctx)

	return NewGenesisState(basePool, terraPool, lastUpdateHeight, params)
}
