package types

import (
	"bytes"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params) GenesisState {

	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}
