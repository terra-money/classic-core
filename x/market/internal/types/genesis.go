package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all market state that must be provided at genesis
type GenesisState struct {
	TerraPoolDelta sdk.Dec `json:"terra_pool_delta" yaml:"terra_pool_delta"`
	Params         Params  `json:"params" yaml:"params"` // market params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(terraPoolDelta sdk.Dec, params Params) GenesisState {
	return GenesisState{
		TerraPoolDelta: terraPoolDelta,
		Params:         params,
	}
}

// DefaultGenesisState returns raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		TerraPoolDelta: sdk.ZeroDec(),
		Params:         DefaultParams(),
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	return data.Params.ValidateBasic()
}

// Equal checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}
