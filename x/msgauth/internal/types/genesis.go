package types

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthorizationEntry hold each authorization information
type AuthorizationEntry struct {
	Granter       sdk.AccAddress `json:"granter" yaml:"granter"`
	Grantee       sdk.AccAddress `json:"grantee" yaml:"grantee"`
	Authorization Authorization  `json:"authorization" yaml:"authorization"`
	Expiration    time.Time      `json:"expiration" yaml:"expiration"`
}

// GenesisState is the struct representation of the export genesis
type GenesisState struct {
	AuthorizationEntries []AuthorizationEntry `json:"authorization_entries" yaml:"authorization_entries"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(entries []AuthorizationEntry) GenesisState {
	return GenesisState{
		AuthorizationEntries: entries,
	}
}

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	return nil
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AuthorizationEntries: []AuthorizationEntry{},
	}
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
