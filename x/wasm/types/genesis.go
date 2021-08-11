package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, lastCodeID, lastInstanceID uint64, codes []Code, contracts []Contract) *GenesisState {
	return &GenesisState{
		Params:         params,
		LastCodeID:     lastCodeID,
		LastInstanceID: lastInstanceID,
		Codes:          codes,
		Contracts:      contracts,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		LastCodeID:     0,
		LastInstanceID: 0,
		Codes:          []Code{},
		Contracts:      []Contract{},
	}
}

// ValidateGenesis performs basic validation of wasm genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data *GenesisState) error {

	if uint64(len(data.Codes)) != data.LastCodeID {
		return sdkerrors.Wrap(ErrInvalidGenesis, "the number of codes is not met with LastCodeID")
	}

	if uint64(len(data.Contracts)) != data.LastInstanceID {
		return sdkerrors.Wrap(ErrInvalidGenesis, "the number of contracts is not met with LastInstanceID")
	}

	return data.Params.Validate()
}

// GetGenesisStateFromAppState returns x/market GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
