package keeper

import (
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetMarketAccount returns market ModuleAccount
func (k Keeper) GetMarketAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetBurnModuleAccount returns burn ModuleAccount
func (k Keeper) GetBurnModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
