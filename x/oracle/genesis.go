package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	FeederDelegations map[string]sdk.AccAddress `json:"feeder_delegations"`
	Params            Params                    `json:"params"` // oracle params
}

// NewGenesisState creates new oracle GenesisState
func NewGenesisState(params Params, feederDelegations map[string]sdk.AccAddress) GenesisState {
	return GenesisState{
		Params:            params,
		FeederDelegations: feederDelegations,
	}
}

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            DefaultParams(),
		FeederDelegations: map[string]sdk.AccAddress{},
	}
}

// InitGenesis creates new oracle genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	for delegator, delegateeAddr := range data.FeederDelegations {
		delegatorAddr, err := sdk.ValAddressFromBech32(delegator)
		if err != nil {
			panic(err)
		}

		keeper.SetFeedDelegate(ctx, delegatorAddr, delegateeAddr)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	feederDelegations := make(map[string]sdk.AccAddress)
	keeper.iterateFeederDelegations(ctx, func(delegatee sdk.AccAddress, delegator sdk.ValAddress) bool {
		feederDelegations[delegator.String()] = delegatee
		return false
	})
	return NewGenesisState(params, feederDelegations)
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	for delegator := range data.FeederDelegations {

		_, err := sdk.ValAddressFromBech32(delegator)
		if err != nil {
			return err
		}

	}

	return validateParams(data.Params)
}
