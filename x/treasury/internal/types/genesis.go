package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all treasury state that must be provided at genesis
type GenesisState struct {
	Params               Params             `json:"params" yaml:"params"` // market params
	TaxRate              sdk.Dec            `json:"tax_rate" yaml:"tax_rate"`
	RewardWeight         sdk.Dec            `json:"reward_weight" yaml:"reward_weight"`
	TaxCaps              map[string]sdk.Int `json:"tax_caps" yaml:"tax_caps"`
	TaxProceed           sdk.Coins          `json:"tax_proceed" yaml:"tax_proceed"`
	EpochInitialIssuance sdk.Coins          `json:"epoch_initial_issuance" yaml:"epoch_initial_issuance"`
	TRs                  []sdk.Dec          `json:"TRs" yaml:"TRs"`
	SRs                  []sdk.Dec          `json:"SRs" yaml:"SRs"`
	TSLs                 []sdk.Int          `json:"TSLs" yaml:"TSLs"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, taxRate sdk.Dec, rewardWeight sdk.Dec,
	taxCaps map[string]sdk.Int, taxProceed sdk.Coins,
	epochInitialIssuance sdk.Coins, TRs []sdk.Dec, SRs []sdk.Dec, TSLs []sdk.Int) GenesisState {
	return GenesisState{
		Params:               params,
		TaxRate:              taxRate,
		RewardWeight:         rewardWeight,
		TaxCaps:              taxCaps,
		TaxProceed:           taxProceed,
		EpochInitialIssuance: epochInitialIssuance,
		TRs:                  TRs,
		SRs:                  SRs,
		TSLs:                 TSLs,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               DefaultParams(),
		TaxRate:              DefaultTaxRate,
		RewardWeight:         DefaultRewardWeight,
		TaxCaps:              make(map[string]sdk.Int),
		TaxProceed:           sdk.Coins{},
		EpochInitialIssuance: sdk.Coins{},
		TRs:                  []sdk.Dec{},
		SRs:                  []sdk.Dec{},
		TSLs:                 []sdk.Int{},
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {

	if data.TaxRate.LT(data.Params.TaxPolicy.RateMin) || data.TaxRate.GT(data.Params.TaxPolicy.RateMax) {
		return fmt.Errorf("tax-rate must less than RateMax(%s) and bigger than RateMin(%s)", data.Params.TaxPolicy.RateMax, data.Params.TaxPolicy.RateMin)
	}

	if data.RewardWeight.LT(data.Params.RewardPolicy.RateMin) || data.RewardWeight.GT(data.Params.RewardPolicy.RateMax) {
		return fmt.Errorf("reward-weight must less than WeightMax(%s) and bigger than RateMin(%s)", data.Params.RewardPolicy.RateMax, data.Params.RewardPolicy.RateMin)
	}

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
