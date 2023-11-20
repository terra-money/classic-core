package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/market/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Replenishes each pools towards equilibrium
	k.ReplenishPools(ctx)
}
