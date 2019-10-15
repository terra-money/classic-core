package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// ParamKeyTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// VotePeriod returns the number of blocks during which voting takes place.
func (k Keeper) VotePeriod(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVotePeriod, &res)
	return
}

// VoteThreshold returns the minimum percentage of votes that must be received for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVoteThreshold, &res)
	return
}

// RewardBand returns the ratio of allowable price error that can be rewared
func (k Keeper) RewardBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardBand, &res)
	return
}

// RewardDistributionPeriod returns the number of blocks of the the period during which seigiornage reward comes in and then is distributed.
func (k Keeper) RewardDistributionPeriod(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardDistributionPeriod, &res)
	return
}

// Whitelist returns the denom list that can be acitivated
func (k Keeper) Whitelist(ctx sdk.Context) (res types.DenomList) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWhitelist, &res)
	return
}

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
