package v05

import (
	v04treasury "github.com/terra-money/core/x/treasury/legacy/v04"
	v05treasury "github.com/terra-money/core/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrate accepts exported v0.4 x/treasury and
// migrates it to v0.5 x/treasury genesis state. The migration includes:
//
// - Merge Epoch genesis data to EpochState from x/treasury genesis state.
// - Update RewardWeight to one to burn all seigniorage
// - Update Params.RewardPolicy so that RewardWeight does not change.
// - Re-encode in v0.5 GenesisState.
func Migrate(
	treasuryGenState v04treasury.GenesisState,
) *v05treasury.GenesisState {
	i := 0
	taxCaps := make([]v05treasury.TaxCap, len(treasuryGenState.TaxCaps))
	for denom, cap := range treasuryGenState.TaxCaps {
		taxCaps[i] = v05treasury.TaxCap{
			Denom:  denom,
			TaxCap: cap,
		}

		i++
	}

	// Remove cumulative height dependencies
	cumulativeEpochs := int(treasuryGenState.CumulativeHeight / int64(v04treasury.BlocksPerWeek))
	epochStates := make([]v05treasury.EpochState, len(treasuryGenState.TRs)-cumulativeEpochs)
	for i := range treasuryGenState.TRs {
		if i < cumulativeEpochs {
			continue
		}

		epochStates[i] = v05treasury.EpochState{
			Epoch:             uint64(i - cumulativeEpochs),
			TaxReward:         treasuryGenState.TRs[i],
			SeigniorageReward: treasuryGenState.SRs[i],
			TotalStakedLuna:   treasuryGenState.TSLs[i],
		}
	}

	return &v05treasury.GenesisState{
		EpochInitialIssuance: treasuryGenState.EpochInitialIssuance,
		EpochStates:          epochStates,
		RewardWeight:         sdk.OneDec(),
		TaxCaps:              taxCaps,
		TaxProceeds:          treasuryGenState.TaxProceed,
		TaxRate:              treasuryGenState.TaxRate,
		Params: v05treasury.Params{
			TaxPolicy: v05treasury.PolicyConstraints{
				RateMin:       treasuryGenState.Params.TaxPolicy.RateMin,
				RateMax:       treasuryGenState.Params.TaxPolicy.RateMax,
				Cap:           treasuryGenState.Params.TaxPolicy.Cap,
				ChangeRateMax: treasuryGenState.Params.TaxPolicy.ChangeRateMax,
			},
			RewardPolicy: v05treasury.PolicyConstraints{
				RateMin:       sdk.ZeroDec(),
				RateMax:       sdk.OneDec(),
				Cap:           treasuryGenState.Params.RewardPolicy.Cap,
				ChangeRateMax: sdk.ZeroDec(),
			},
			MiningIncrement:         treasuryGenState.Params.MiningIncrement,
			SeigniorageBurdenTarget: treasuryGenState.Params.SeigniorageBurdenTarget,
			WindowShort:             uint64(treasuryGenState.Params.WindowShort),
			WindowLong:              uint64(treasuryGenState.Params.WindowLong),
			WindowProbation:         uint64(treasuryGenState.Params.WindowProbation),
		},
	}
}
