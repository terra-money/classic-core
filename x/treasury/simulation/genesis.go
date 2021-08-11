package simulation

//DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/types"
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
func GenWindowShort(r *rand.Rand) uint64 {
	return uint64(1 + r.Intn(12))
}

// GenWindowLong randomized WindowLong
func GenWindowLong(r *rand.Rand) uint64 {
	return uint64(12 + r.Intn(24))
}

// GenWindowProbation randomized WindowProbation
func GenWindowProbation(r *rand.Rand) uint64 {
	return uint64(1 + r.Intn(6))
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

	var windowShort uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, windowShortKey, &windowShort, simState.Rand,
		func(r *rand.Rand) { windowShort = GenWindowShort(r) },
	)

	var windowLong uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, windowLongKey, &windowLong, simState.Rand,
		func(r *rand.Rand) { windowLong = GenWindowLong(r) },
	)

	var windowProbation uint64
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
		[]types.TaxCap{},
		sdk.Coins{},
		sdk.Coins{},
		[]types.EpochState{},
	)

	bz, err := json.MarshalIndent(&treasuryGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected randomly generated market parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(treasuryGenesis)
}
