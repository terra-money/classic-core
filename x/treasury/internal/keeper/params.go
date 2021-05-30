package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/treasury/internal/types"
)

// TaxPolicy defines constraints for TaxRate
func (k Keeper) TaxPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyTaxPolicy, &res)
	return
}

// RewardPolicy defines constraints for RewardWeight
func (k Keeper) RewardPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardPolicy, &res)
	return
}

// SeigniorageBurdenTarget defines fixed target for the Seigniorage Burden. Between 0 and 1.
func (k Keeper) SeigniorageBurdenTarget(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySeigniorageBurdenTarget, &res)
	return
}

// MiningIncrement is a factor used to determine how fast MRL should grow over time
func (k Keeper) MiningIncrement(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMiningIncrement, &res)
	return
}

// WindowShort is a short period window for moving average
func (k Keeper) WindowShort(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowShort, &res)
	return
}

// WindowLong is a long period window for moving average
func (k Keeper) WindowLong(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowLong, &res)
	return
}

// WindowProbation is a period of time to prevent updates
func (k Keeper) WindowProbation(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWindowProbation, &res)
	return
}

// GetParams returns the total set of treasury parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of treasury parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
