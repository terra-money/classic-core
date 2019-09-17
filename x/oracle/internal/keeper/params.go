package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// VotePeriod
func (k Keeper) VotePeriod(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVotePeriod, &res)
	return
}

// VoteThreshold
func (k Keeper) VoteThreshold(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVoteThreshold, &res)
	return
}

// RewardBand
func (k Keeper) RewardBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardBand, &res)
	return
}

// RewardFraction
func (k Keeper) RewardFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardFraction, &res)
	return
}

// VotesWindow
func (k Keeper) VotesWindow(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyVotesWindow, &res)
	return
}

// MinValidVotesPerWindow
func (k Keeper) MinValidVotesPerWindow(ctx sdk.Context) (res int64) {
	var minValidVotesPerWindow sdk.Dec
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinValidVotesPerWindow, &minValidVotesPerWindow)
	signedBlocksWindow := k.VotesWindow(ctx)

	// NOTE: RoundInt64 will never panic as minValidVotesPerWindow is less than 1.
	return minValidVotesPerWindow.MulInt64(signedBlocksWindow).RoundInt64()
}

// SlashFraction
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySlashFraction, &res)
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
