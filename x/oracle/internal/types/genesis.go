package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params            Params                    `json:"params" yaml:"params"`
	FeederDelegations map[string]sdk.AccAddress `json:"feeder_delegations" yaml:"feeder_delegations"`
	Prices            map[string]sdk.Dec        `json:"prices" yaml:"prices"`
	PricePrevotes     []PricePrevote            `json:"price_prevotes" yaml:"price_prevotes"`
	PriceVotes        []PriceVote               `json:"price_votes" yaml:"price_votes"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, pricePrevotes []PricePrevote,
	priceVotes []PriceVote, prices map[string]sdk.Dec,
	feederDelegations map[string]sdk.AccAddress,
) GenesisState {

	return GenesisState{
		Params:            params,
		PricePrevotes:     pricePrevotes,
		PriceVotes:        priceVotes,
		Prices:            prices,
		FeederDelegations: feederDelegations,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            DefaultParams(),
		PricePrevotes:     []PricePrevote{},
		PriceVotes:        []PriceVote{},
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
