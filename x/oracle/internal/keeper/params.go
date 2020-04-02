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

// RewardBand returns the ratio of allowable exchange rate error that a validator can be rewared
func (k Keeper) RewardBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardBand, &res)
	return
}

// RewardDistributionWindow returns the number of vote periods during which seigiornage reward comes in and then is distributed.
func (k Keeper) RewardDistributionWindow(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardDistributionWindow, &res)
	return
}

// Whitelist returns the denom list that can be activated
func (k Keeper) Whitelist(ctx sdk.Context) (res types.DenomList) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWhitelist, &res)
	return
}

// SlashFraction returns oracle voting penalty rate
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySlashFraction, &res)
	return
}

// SlashWindow returns # of vote period for oracle slashing
func (k Keeper) SlashWindow(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.ParamStoreKeySlashWindow, &res)
	return
}

// MinValidPerWindow returns oracle slashing threshold
func (k Keeper) MinValidPerWindow(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMinValidPerWindow, &res)
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
