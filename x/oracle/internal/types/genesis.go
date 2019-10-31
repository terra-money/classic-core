package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params            Params                    `json:"params" yaml:"params"`
	FeederDelegations map[string]sdk.AccAddress `json:"feeder_delegations" yaml:"feeder_delegations"`
	Prices            map[string]sdk.Dec        `json:"exchangeRates" yaml:"exchangeRates"`
	Prevotes          []Prevote                 `json:"exchangeRate_prevotes" yaml:"exchangeRate_prevotes"`
	Votes             []Vote                    `json:"exchangeRate_votes" yaml:"exchangeRate_votes"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, Prevotes []Prevote,
	Votes []Vote, exchangeRates map[string]sdk.Dec,
	feederDelegations map[string]sdk.AccAddress,
) GenesisState {

	return GenesisState{
		Params:            params,
		Prevotes:          Prevotes,
		Votes:             Votes,
		Prices:            exchangeRates,
		FeederDelegations: feederDelegations,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            DefaultParams(),
		Prevotes:          []Prevote{},
		Votes:             []Vote{},
		Prices:            make(map[string]sdk.Dec),
		FeederDelegations: make(map[string]sdk.AccAddress),
	}
}

// ValidateGenesis validates the oracle genesis parameters
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
