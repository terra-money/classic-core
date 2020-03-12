package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
	"time"
)

// BidPeriod returns the time period for Luna holders to bid on an auction
func (k Keeper) BidPeriod(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyBidPeriod, &res)
	return
}

// RevealPeriod returns the time period for Luna holders to reveal the proof of a bid on an auction
func (k Keeper) RevealPeriod(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRevealPeriod, &res)
	return
}

// GracePeriod returns the time period for name registry owners to extend the ownership
func (k Keeper) GracePeriod(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyGracePeriod, &res)
	return
}

// RenewalInterval returns the time frequency of renewal of ownership.
func (k Keeper) RenewalInterval(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRenewalInterval, &res)
	return
}

// MinDeposit returns minimum deposit required to bid
func (k Keeper) MinDeposit(ctx sdk.Context) (res sdk.Coin) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinDeposit, &res)
	return
}

// RootName returns enforced first level name of the nameservice module
func (k Keeper) RootName(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRootName, &res)
	return
}

// RootName returns minimum name length for auction
func (k Keeper) MinNameLength(ctx sdk.Context) (res int) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinNameLength, &res)
	return
}

// RenewalFees returns required renewal fees to renew name registry during renewal interval
func (k Keeper) RenewalFees(ctx sdk.Context) (res types.RenewalFees) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRenewalFees, &res)
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
