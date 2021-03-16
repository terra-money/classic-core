package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

var _ codectypes.UnpackInterfacesMessage = AuthorizationEntry{}
var _ codectypes.UnpackInterfacesMessage = GenesisState{}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(entries []AuthorizationEntry) *GenesisState {
	return &GenesisState{
		AuthorizationEntries: entries,
	}
}

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data *GenesisState) error {
	return nil
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AuthorizationEntries: []AuthorizationEntry{},
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (entry AuthorizationEntry) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization AuthorizationI
	return unpacker.UnpackAny(entry.Authorization, &authorization)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, a := range data.AuthorizationEntries {
		err := a.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAuthorization returns the grant Authorization
func (entry AuthorizationEntry) GetAuthorization() AuthorizationI {
	authorization, ok := entry.Authorization.GetCachedValue().(AuthorizationI)
	if !ok {
		return nil
	}
	return authorization
}
