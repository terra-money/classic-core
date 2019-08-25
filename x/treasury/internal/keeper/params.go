package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/treasury/internal/types"
)

// ParamTable for treasury module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// TaxPolicy
func (k Keeper) TaxPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyTaxPolicy, &res)
	return
}

// RewardPolicy
func (k Keeper) RewardPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardPolicy, &res)
	return
}

// SeigniorageBurdenTarget
func (k Keeper) SeigniorageBurdenTarget(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySeigniorageBurdenTarget, &res)
	return
}

// MiningIncrement
func (k Keeper) MiningIncrement(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMiningIncrement, &res)
	return
}

// WindowShort
func (k Keeper) WindowShort(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowShort, &res)
	return
}

// WindowLong
func (k Keeper) WindowLong(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowLong, &res)
	return
}

// WindowProbation
func (k Keeper) WindowProbation(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowProbation, &res)
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
