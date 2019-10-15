package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// AfterValidatorBonded register voting info for a validator
func (k Keeper) AfterValidatorBonded(ctx sdk.Context, _ sdk.ConsAddress, address sdk.ValAddress) {
	// Update the signing info start height or create a new signing info
	_, found := k.getVotingInfo(ctx, address)
	if !found {
		votingInfo := types.NewVotingInfo(
			address,
			ctx.BlockHeight(),
			0,
			0,
		)
		k.SetVotingInfo(ctx, address, votingInfo)
	}
}

// Hooks wrapper struct for oracle keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Keeper{}

// Return the wrapper struct
func (k Keeper) StakingHooks() Hooks {
	return Hooks{k}
}

// Implements StakingHooks
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.AfterValidatorBonded(ctx, consAddr, valAddr)
}

// nolint - unused hooks
func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                            {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)  {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}
