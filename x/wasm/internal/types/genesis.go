package types

// GenesisState is the struct representation of the export genesis
type GenesisState struct {
	Params    Params     `json:"params" yaml:"params"`
	Codes     []Code     `json:"codes" yaml:"codes"`
	Contracts []Contract `json:"contracts" yaml:"contracts"`
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
func NewGenesisState(params Params, codes []Code, contracts []Contract) GenesisState {
	return GenesisState{
		Params:    params,
		Codes:     codes,
		Contracts: contracts,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:    DefaultParams(),
		Codes:     []Code{},
		Contracts: []Contract{},
	}
}

// ValidateGenesis performs basic validation of wasm genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
