package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/terra-money/core/x/market/internal/types"
)

// GetMarketAccount returns market ModuleAccount
func (k Keeper) GetMarketAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
