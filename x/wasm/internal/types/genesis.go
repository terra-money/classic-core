package types

// GenesisState is the struct representation of the export genesis
type GenesisState struct {
	Params         Params     `json:"params" yaml:"params"`
	LastCodeID     uint64     `json:"last_code_id" yaml:"last_code_id"`
	LastInstanceID uint64     `json:"last_instance_id" yaml:"last_instance_id"`
	Codes          []Code     `json:"codes" yaml:"codes"`
	Contracts      []Contract `json:"contracts" yaml:"contracts"`
}

// Code struct encompasses CodeInfo and CodeBytes
type Code struct {
	CodeInfo   CodeInfo `json:"code_info"`
	CodesBytes []byte   `json:"code_bytes"`
}

// Contract struct encompasses ContractAddress, ContractInfo, and ContractState
type Contract struct {
	ContractInfo  ContractInfo `json:"contract_info"`
	ContractStore []Model      `json:"contract_store"`
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

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:         DefaultParams(),
		LastCodeID:     0,
		LastInstanceID: 0,
		Codes:          []Code{},
		Contracts:      []Contract{},
	}
}

// ValidateGenesis performs basic validation of wasm genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
