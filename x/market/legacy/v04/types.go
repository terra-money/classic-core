// DONTCOVER
// nolint
package v04

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName nolint
	ModuleName = "market"
)

type (
	// Params market parameters
	Params struct {
		BasePool           sdk.Dec `json:"base_pool" yaml:"base_pool"`
		PoolRecoveryPeriod int64   `json:"pool_recovery_period" yaml:"pool_recovery_period"`
		MinStabilitySpread sdk.Dec `json:"min_spread" yaml:"min_spread"`
	}

	// GenesisState is the struct representation of the export genesis
	GenesisState struct {
		TerraPoolDelta sdk.Dec `json:"terra_pool_delta" yaml:"terra_pool_delta"`
		Params         Params  `json:"params" yaml:"params"` // market params
	}
)
