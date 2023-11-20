package keeper

import (
	"github.com/classic-terra/core/v2/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TaxPolicy defines constraints for TaxRate
func (k Keeper) TaxPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.KeyTaxPolicy, &res)
	return
}

// RewardPolicy defines constraints for RewardWeight
func (k Keeper) RewardPolicy(ctx sdk.Context) (res types.PolicyConstraints) {
	k.paramSpace.Get(ctx, types.KeyRewardPolicy, &res)
	return
}

// SeigniorageBurdenTarget defines fixed target for the Seigniorage Burden. Between 0 and 1.
func (k Keeper) SeigniorageBurdenTarget(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySeigniorageBurdenTarget, &res)
	return
}

// MiningIncrement is a factor used to determine how fast MRL should grow over time
func (k Keeper) MiningIncrement(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMiningIncrement, &res)
	return
}

// WindowShort is a short period window for moving average
func (k Keeper) WindowShort(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyWindowShort, &res)
	return
}

// WindowLong is a long period window for moving average
func (k Keeper) WindowLong(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyWindowLong, &res)
	return
}

// WindowProbation is a period of time to prevent updates
func (k Keeper) WindowProbation(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyWindowProbation, &res)
	return
}

func (k Keeper) GetBurnSplitRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyBurnTaxSplit, &res)
	return
}

func (k Keeper) SetBurnSplitRate(ctx sdk.Context, burnTaxSplit sdk.Dec) {
	k.paramSpace.Set(ctx, types.KeyBurnTaxSplit, burnTaxSplit)
}

func (k Keeper) GetMinInitialDepositRatio(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMinInitialDepositRatio, &res)
	return
}

func (k Keeper) SetMinInitialDepositRatio(ctx sdk.Context, minInitialDepositRatio sdk.Dec) {
	k.paramSpace.Set(ctx, types.KeyMinInitialDepositRatio, minInitialDepositRatio)
}

// GetParams returns the total set of treasury parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSetIfExists(ctx, &params)
	return params
}

// SetParams sets the total set of treasury parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
