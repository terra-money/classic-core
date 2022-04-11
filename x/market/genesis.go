package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetTerraPoolDelta(ctx, data.TerraPoolDelta)

	// check if the module account exists
	moduleAcc := keeper.GetMarketAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// check if the burn module account exists
	burnModuleAcc := keeper.GetBurnModuleAccount(ctx)
	if burnModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BurnModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)
	terraPoolDelta := keeper.GetTerraPoolDelta(ctx)

	return types.NewGenesisState(terraPoolDelta, params)
}
