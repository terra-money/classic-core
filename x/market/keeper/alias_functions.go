package keeper

import (
	"github.com/classic-terra/core/v2/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetMarketAccount returns market ModuleAccount
func (k Keeper) GetMarketAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
