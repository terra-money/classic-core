package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/market/internal/types"
)

// ParamTable for market module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// DailyTerraLiquidityRatio is the ratio of Terra's market cap which will be made available per day
func (k Keeper) DailyTerraLiquidityRatio(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyDailyTerraLiquidityRatio, &res)
	return
}

// MinSpread is the minimum swap fee(spread)
func (k Keeper) MinSpread(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinSpread, &res)
	return
}

// PoolUpdateInterval is the swap rate of base denom
func (k Keeper) PoolUpdateInterval(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolUpdateInterval, &res)
	return
}

// TobinTax is a tax on all spot conversions of one TERRA into another TERRA
func (k Keeper) TobinTax(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParmamStoreKeyTobinTax, &res)
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
