package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, rates []ValidatorCommissionRate) *GenesisState {
	return &GenesisState{
		Params:                   params,
		ValidatorCommissionRates: rates,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	emptySet := []ValidatorCommissionRate{}
	return &GenesisState{
		Params:                   DefaultParams(),
		ValidatorCommissionRates: emptySet,
	}
}
