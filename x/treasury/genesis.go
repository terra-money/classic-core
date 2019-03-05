package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all treasury state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // treasury params
}

func NewGenesisState(params Params /*, genesisIssuance map[string]sdk.Int*/) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// new oracle genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	return NewGenesisState(params)
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	return validateParams(data.Params)
}
