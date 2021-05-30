package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/terra-money/core/x/treasury/internal/types"
)

// GetTreasuryAccount returns treasury ModuleAccount
func (k Keeper) GetTreasuryAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
