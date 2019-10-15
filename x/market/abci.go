package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {

	// Replenishes each pools towards equilibrium
	k.ReplenishPools(ctx)

}
