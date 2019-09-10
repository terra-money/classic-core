package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {

	// Replenishes each pools towards equilibrium
	k.ReplenishPools(ctx)

	// Update pools at the last block of every interval
	// Retry update when inactive state
	if !core.IsPeriodLastBlock(ctx, k.PoolUpdateInterval(ctx)) && k.IsMarketActive(ctx) {
		return
	}

	basePool, err := k.UpdatePools(ctx)

	if err != nil {
		// TODO - check log level
		k.Logger(ctx).Error(fmt.Sprintf("Failed to update BasePool: %s", err))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventPoolUpdate,
			sdk.NewAttribute(types.AttributeKeyBasePool, basePool.String()),
		),
	)
}
