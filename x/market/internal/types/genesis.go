package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all market state that must be provided at genesis
type GenesisState struct {
	BasePool         sdk.Dec `json:"base_pool" yaml:"base_pool"`
	TerraPool        sdk.Dec `json:"terra_pool" yaml:"terra_pool"`
	LastUpdateHeight int64   `json:"last_update_height" yaml:"last_update_height"`
	Params           Params  `json:"params" yaml:"params"` // market params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(basePool, terraPool sdk.Dec, lastUpdateHeight int64, params Params) GenesisState {
	return GenesisState{
		BasePool:         basePool,
		TerraPool:        terraPool,
		LastUpdateHeight: lastUpdateHeight,
		Params:           params,
	}
}

// DefaultGenesisState returns raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		BasePool:         sdk.ZeroDec(),
		TerraPool:        sdk.ZeroDec(),
		LastUpdateHeight: 0,
		Params:           DefaultParams(),
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
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
