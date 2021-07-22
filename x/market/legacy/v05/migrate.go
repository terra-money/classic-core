package v05

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v04market "github.com/terra-money/core/x/market/legacy/v04"
	v05market "github.com/terra-money/core/x/market/types"
)

// Migrate accepts exported v0.4 x/market and
// migrates it to v0.5 x/market genesis state. The migration includes:
//
// - Split BasePool to MintPool and Burn Pool from x/market genesis state.
// - Re-encode in v0.5 GenesisState.
func Migrate(
	marketGenState v04market.GenesisState,
) *v05market.GenesisState {
	return &v05market.GenesisState{
		TerraPoolDelta: sdk.ZeroDec(),
		Params: v05market.Params{
			BasePool:           marketGenState.Params.BasePool,
			PoolRecoveryPeriod: uint64(marketGenState.Params.PoolRecoveryPeriod),
			MinStabilitySpread: marketGenState.Params.MinStabilitySpread,
		},
	}
}
