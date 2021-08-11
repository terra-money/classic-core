package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/terra-money/core/x/treasury/types"
)

// GetTreasuryModuleAccount returns treasury ModuleAccount
func (k Keeper) GetTreasuryModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetBurnModuleAccount returns burn ModuleAccount
func (k Keeper) GetBurnModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.BurnModuleName)
}
