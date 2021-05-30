package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/market/keeper"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	// Replenishes each pools towards equilibrium
	k.ReplenishPools(ctx)

}
