package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/keeper"
)

// BeginBlocker is called at the begin of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

	// Make min spread to one to disable swap
	if ctx.ChainID() == core.ColumbusChainID && ctx.BlockHeight() == core.SwapDisableForkHeight {
		params := k.GetParams(ctx)
		params.MinStabilitySpread = sdk.OneDec()
		k.SetParams(ctx, params)
	}

}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	// Replenishes each pools towards equilibrium
	k.ReplenishPools(ctx)

}
