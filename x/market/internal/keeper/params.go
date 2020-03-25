package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/market/internal/types"
)

// ParamKeyTable for market module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// BasePool is Terra liquidity pool(usdr unit) which will be made available per PoolRecoveryPeriod
func (k Keeper) BasePool(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyBasePool, &res)
	return
}

// MinStabilitySpread is the minimum spread applied to swaps to / from Luna.
// Intended to prevent swing trades exploiting oracle period delays
func (k Keeper) MinStabilitySpread(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinStabilitySpread, &res)
	return
}

// PoolRecoveryPeriod is the period required to recover Terra&Luna Pools to the BasePool
func (k Keeper) PoolRecoveryPeriod(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolRecoveryPeriod, &res)
	return
}

// GetParams returns the total set of market parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of market parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
