package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/internal/types"
)

// Simulation parameter constants
const (
	taxPolicyKey               = "tax_policy"
	rewardPolicyKey            = "reward_policy"
	seigniorageBurdenTargetKey = "seigniorage_burden_target"
	miningIncrementKey         = "mining_increment"
	windowShortKey             = "window_short"
	windowLongKey              = "window_long"
	windowProbationKey         = "window_probation"
)

// GenTaxPolicy randomized TaxPolicy
func GenTaxPolicy(r *rand.Rand) types.PolicyConstraints {
	return types.PolicyConstraints{
		RateMin:       sdk.NewDecWithPrec(int64(r.Intn(5)+1), 3),
		RateMax:       sdk.NewDecWithPrec(6, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(5)+1), 3)),
		Cap:           sdk.NewInt64Coin(core.MicroSDRDenom, 1000000),
		ChangeRateMax: sdk.NewDecWithPrec(25, 5).Add(sdk.NewDecWithPrec(int64(r.Intn(75)), 5)),
	}
}

// GenRewardPolicy randomized RewardPolicy
func GenRewardPolicy(r *rand.Rand) types.PolicyConstraints {
	return types.PolicyConstraints{
		RateMin:       sdk.NewDecWithPrec(int64(r.Intn(5)+1), 3),
		RateMax:       sdk.NewDecWithPrec(6, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(5)+1), 3)),
		Cap:           sdk.NewCoin("unused", sdk.ZeroInt()),
		ChangeRateMax: sdk.NewDecWithPrec(25, 5).Add(sdk.NewDecWithPrec(int64(r.Intn(75)), 5)),
	}
}

// GenSeigniorageBurdenTarget randomized SeigniorageBurdenTarget
func GenSeigniorageBurdenTarget(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(100)), 2)
}

// GenMiningIncrement randomized MiningIncrement
func GenMiningIncrement(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(100+r.Intn(30)), 2)
}

// GenWindowShort randomized WindowShort
func GenWindowShort(r *rand.Rand) int64 {
	return int64(1 + r.Intn(12))
}

// GenWindowLong randomized WindowLong
func GenWindowLong(r *rand.Rand) int64 {
	return int64(12 + r.Intn(24))
}

// GenWindowProbation randomized WindowProbation
func GenWindowProbation(r *rand.Rand) int64 {
	return int64(1 + r.Intn(6))
}

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {

	var taxPolicy types.PolicyConstraints
	simState.AppParams.GetOrGenerate(
		simState.Cdc, taxPolicyKey, &taxPolicy, simState.Rand,
		func(r *rand.Rand) { taxPolicy = GenTaxPolicy(r) },
	)

	var rewardPolicy types.PolicyConstraints
	simState.AppParams.GetOrGenerate(
		simState.Cdc, rewardPolicyKey, &rewardPolicy, simState.Rand,
		func(r *rand.Rand) { rewardPolicy = GenRewardPolicy(r) },
	)

	var seigniorageBurdenTarget sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, seigniorageBurdenTargetKey, &seigniorageBurdenTarget, simState.Rand,
		func(r *rand.Rand) { seigniorageBurdenTarget = GenSeigniorageBurdenTarget(r) },
	)

	var miningIncrement sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, miningIncrementKey, &miningIncrement, simState.Rand,
		func(r *rand.Rand) { miningIncrement = GenMiningIncrement(r) },
	)

	var windowShort int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, windowShortKey, &windowShort, simState.Rand,
		func(r *rand.Rand) { windowShort = GenWindowShort(r) },
	)

	var windowLong int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, windowLongKey, &windowLong, simState.Rand,
		func(r *rand.Rand) { windowLong = GenWindowLong(r) },
	)

	var windowProbation int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, windowProbationKey, &windowProbation, simState.Rand,
		func(r *rand.Rand) { windowProbation = GenWindowProbation(r) },
	)

	treasuryGenesis := types.NewGenesisState(
		types.Params{
			TaxPolicy:               taxPolicy,
			RewardPolicy:            rewardPolicy,
			SeigniorageBurdenTarget: seigniorageBurdenTarget,
			MiningIncrement:         miningIncrement,
			WindowShort:             windowShort,
			WindowLong:              windowLong,
			WindowProbation:         windowProbation,
		},
		taxPolicy.RateMin,
		rewardPolicy.RateMin,
		map[string]sdk.Int{},
		sdk.Coins{},
		sdk.Coins{},
		0,
		[]sdk.Dec{},
		[]sdk.Dec{},
		[]sdk.Int{},
	)

	fmt.Printf("Selected randomly generated treasury parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, treasuryGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(treasuryGenesis)
}
