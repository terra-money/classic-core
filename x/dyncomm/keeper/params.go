package keeper

import (
	"github.com/classic-terra/core/v2/x/dyncomm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetMaxZero(ctx sdk.Context) (ret sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMaxZero, &ret)
	return ret
}

func (k Keeper) GetSlopeBase(ctx sdk.Context) (ret sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySlopeBase, &ret)
	return ret
}

func (k Keeper) GetSlopeVpImpact(ctx sdk.Context) (ret sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySlopeVpImpact, &ret)
	return ret
}

func (k Keeper) GetCap(ctx sdk.Context) (ret sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyCap, &ret)
	return ret
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
