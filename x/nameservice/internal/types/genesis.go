package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params     Params                    `json:"params" yaml:"params"`
	Auctions   map[string]Auction        `json:"auctions" yaml:"auctions"`
	Bids       map[string]Bid            `json:"bids" yaml:"bids"`
	Registries map[string]Registry       `json:"registries" yaml:"registries"`
	Resolves   map[string]sdk.AccAddress `json:"resolves" yaml:"resolves"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, auctions map[string]Auction, bids map[string]Bid,
	registries map[string]Registry, resolves map[string]sdk.AccAddress,
) GenesisState {
	return GenesisState{
		Params:     params,
		Registries: registries,
		Resolves:   resolves,
		Auctions:   auctions,
		Bids:       bids,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:     DefaultParams(),
		Registries: make(map[string]Registry),
		Resolves:   make(map[string]sdk.AccAddress),
		Auctions:   make(map[string]Auction),
		Bids:       make(map[string]Bid),
	}
}

// ValidateGenesis validates the nameservice genesis parameters
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}

// Equal checks whether 2 GenesisState struct are equivalent.
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
