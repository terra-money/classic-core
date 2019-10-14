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

// VotesWindow returns the number of block units on which the penalty is based
func (k Keeper) VotesWindow(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVotesWindow, &res)
	return
}

// MinValidVotesPerWindow returns the minimum number of blocks to avoid slashing in a window
func (k Keeper) MinValidVotesPerWindow(ctx sdk.Context) (res int64) {
	var minValidVotesPerWindow sdk.Dec
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinValidVotesPerWindow, &minValidVotesPerWindow)
	signedBlocksWindow := k.VotesWindow(ctx)

	// NOTE: RoundInt64 will never panic as minValidVotesPerWindow is less than 1.
	return minValidVotesPerWindow.MulInt64(signedBlocksWindow).RoundInt64()
}

// SlashFraction returns the slashing ratio on the delegated token
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySlashFraction, &res)
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
