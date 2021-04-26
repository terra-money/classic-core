// DONTCOVER
// nolint
package v04

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "wasm"
)

type (
	// Params wasm parameters
	Params struct {
		MaxContractSize    uint64 `json:"max_contract_size" yaml:"max_contract_size"`         // allowed max contract bytes size
		MaxContractGas     uint64 `json:"max_contract_gas" yaml:"max_contract_gas"`           // allowed max gas usages per each contract execution
		MaxContractMsgSize uint64 `json:"max_contract_msg_size" yaml:"max_contract_msg_size"` // allowed max contract exe msg bytes size
	}

	// GenesisState is the struct representation of the export genesis
	GenesisState struct {
		Params         Params     `json:"params" yaml:"params"`
		LastCodeID     uint64     `json:"last_code_id" yaml:"last_code_id"`
		LastInstanceID uint64     `json:"last_instance_id" yaml:"last_instance_id"`
		Codes          []Code     `json:"codes" yaml:"codes"`
		Contracts      []Contract `json:"contracts" yaml:"contracts"`
	}
)

// Code struct encompasses CodeInfo and CodeBytes
type Code struct {
	CodeInfo   CodeInfo `json:"code_info"`
	CodesBytes []byte   `json:"code_bytes"`
}

// CodeInfo is data for the uploaded contract WASM code
type CodeInfo struct {
	CodeID   uint64         `json:"code_id"`
	CodeHash []byte         `json:"code_hash"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Contract struct encompasses ContractAddress, ContractInfo, and ContractState
type Contract struct {
	ContractInfo  ContractInfo `json:"contract_info"`
	ContractStore []Model      `json:"contract_store"`
}

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	Address    sdk.AccAddress `json:"address"`
	Owner      sdk.AccAddress `json:"owner"`
	CodeID     uint64         `json:"code_id"`
	InitMsg    []byte         `json:"init_msg"`
	Migratable bool           `json:"migratable"`
}

// Model is a struct that holds a KV pair
type Model struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, lastCodeID, lastInstanceID uint64, codes []Code, contracts []Contract) GenesisState {
	return GenesisState{
		Params:         params,
		LastCodeID:     lastCodeID,
		LastInstanceID: lastInstanceID,
		Codes:          codes,
		Contracts:      contracts,
	}
}
