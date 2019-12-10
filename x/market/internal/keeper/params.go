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

// MinSpread is the minimum swap fee(spread)
func (k Keeper) MinSpread(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinSpread, &res)
	return
}

// PoolRecoveryPeriod is the period required to recover Terra&Luna Pools to the BasePool
func (k Keeper) PoolRecoveryPeriod(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolRecoveryPeriod, &res)
	return
}

// TobinTax is a tax rate on all spot conversions of one Terra into another Terra
func (k Keeper) TobinTax(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParmaStoreKeyTobinTax, &res)
	return
}

// IlliquidTobinTaxList is the exceptions that have to pay a higher tobin tax due to illiquidity
// TobinTax will be used for the denoms which are not in the list
func (k Keeper) IlliquidTobinTaxList(ctx sdk.Context) (res types.TobinTaxList) {
	k.paramSpace.Get(ctx, types.ParmaStoreKeyIlliquidTobinTaxList, &res)
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
