package keeper

import (
	"github.com/classic-terra/core/v2/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BasePool is liquidity pool(usdr unit) which will be made available per PoolRecoveryPeriod
func (k Keeper) BasePool(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyBasePool, &res)
	return
}

// MinStabilitySpread is the minimum spread applied to swaps to / from Luna.
// Intended to prevent swing trades exploiting oracle period delays
func (k Keeper) MinStabilitySpread(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMinStabilitySpread, &res)
	return
}

// PoolRecoveryPeriod is the period required to recover Terra&Luna Pools to the MintBasePool & BurnBasePool
func (k Keeper) PoolRecoveryPeriod(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyPoolRecoveryPeriod, &res)
	return
}

// GetParams returns the total set of market parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSetIfExists(ctx, &params)
	return params
}

// SetParams sets the total set of market parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
