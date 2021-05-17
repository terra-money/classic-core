package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, taxRate sdk.Dec, rewardWeight sdk.Dec,
	taxCaps []TaxCap, taxProceeds sdk.Coins, epochInitialIssuance sdk.Coins,
	epochStates []EpochState) *GenesisState {
	return &GenesisState{
		Params:               params,
		TaxRate:              taxRate,
		RewardWeight:         rewardWeight,
		TaxCaps:              taxCaps,
		TaxProceeds:          taxProceeds,
		EpochInitialIssuance: epochInitialIssuance,
		EpochStates:          epochStates,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		TaxRate:              DefaultTaxRate,
		RewardWeight:         DefaultRewardWeight,
		TaxCaps:              []TaxCap{},
		TaxProceeds:          sdk.Coins{},
		EpochInitialIssuance: sdk.Coins{},
		EpochStates:          []EpochState{},
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data *GenesisState) error {

	if data.TaxRate.LT(data.Params.TaxPolicy.RateMin) || data.TaxRate.GT(data.Params.TaxPolicy.RateMax) {
		return fmt.Errorf("tax_rate must less than RateMax(%s) and bigger than RateMin(%s)", data.Params.TaxPolicy.RateMax, data.Params.TaxPolicy.RateMin)
	}

	if data.RewardWeight.LT(data.Params.RewardPolicy.RateMin) || data.RewardWeight.GT(data.Params.RewardPolicy.RateMax) {
		return fmt.Errorf("reward_weight must less than WeightMax(%s) and bigger than RateMin(%s)", data.Params.RewardPolicy.RateMax, data.Params.RewardPolicy.RateMin)
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
