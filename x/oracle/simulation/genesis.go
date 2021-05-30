package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/internal/types"
)

// Simulation parameter constants
const (
	votePeriodKey               = "vote_period"
	voteThresholdKey            = "vote_threshold"
	rewardBandKey               = "reward_band"
	rewardDistributionWindowKey = "reward_distribution_window"
	slashFractionKey            = "slash_fraction"
	slashWindowKey              = "slash_window"
	minValidPerWindowKey        = "min_valid_per_window"
)

// GenVotePeriod randomized VotePeriod
func GenVotePeriod(r *rand.Rand) int64 {
	return int64(1 + r.Intn(100))
}

// GenVoteThreshold randomized VoteThreshold
func GenVoteThreshold(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(333, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(333)), 3))
}

// GenRewardBand randomized RewardBand
func GenRewardBand(r *rand.Rand) sdk.Dec {
	return sdk.ZeroDec().Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenRewardDistributionWindow randomized RewardDistributionWindow
func GenRewardDistributionWindow(r *rand.Rand) int64 {
	return int64(100 + r.Intn(100000))
}

// GenSlashFraction randomized SlashFraction
func GenSlashFraction(r *rand.Rand) sdk.Dec {
	return sdk.ZeroDec().Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenSlashWindow randomized SlashWindow
func GenSlashWindow(r *rand.Rand) int64 {
	return int64(100 + r.Intn(100000))
}

// GenMinValidPerWindow randomized MinValidPerWindow
func GenMinValidPerWindow(r *rand.Rand) sdk.Dec {
	return sdk.ZeroDec().Add(sdk.NewDecWithPrec(int64(r.Intn(500)), 3))
}

// RandomizedGenState generates a random GenesisState for oracle
func RandomizedGenState(simState *module.SimulationState) {

	var votePeriod int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, votePeriodKey, &votePeriod, simState.Rand,
		func(r *rand.Rand) { votePeriod = GenVotePeriod(r) },
	)

	var voteThreshold sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, voteThresholdKey, &voteThreshold, simState.Rand,
		func(r *rand.Rand) { voteThreshold = GenVoteThreshold(r) },
	)

	var rewardBand sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, rewardBandKey, &rewardBand, simState.Rand,
		func(r *rand.Rand) { rewardBand = GenRewardBand(r) },
	)

	var rewardDistributionWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, rewardDistributionWindowKey, &rewardDistributionWindow, simState.Rand,
		func(r *rand.Rand) { rewardDistributionWindow = GenRewardDistributionWindow(r) },
	)

	var slashFraction sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, slashFractionKey, &slashFraction, simState.Rand,
		func(r *rand.Rand) { slashFraction = GenSlashFraction(r) },
	)

	var slashWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, slashWindowKey, &slashWindow, simState.Rand,
		func(r *rand.Rand) { slashWindow = GenSlashWindow(r) },
	)

	var minValidPerWindow sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, minValidPerWindowKey, &minValidPerWindow, simState.Rand,
		func(r *rand.Rand) { minValidPerWindow = GenMinValidPerWindow(r) },
	)

	oracleGenesis := types.NewGenesisState(
		types.Params{
			VotePeriod:               votePeriod,
			VoteThreshold:            voteThreshold,
			RewardBand:               rewardBand,
			RewardDistributionWindow: rewardDistributionWindow,
			Whitelist: types.DenomList{
				{core.MicroKRWDenom, types.DefaultTobinTax},
				{core.MicroSDRDenom, types.DefaultTobinTax},
				{core.MicroUSDDenom, types.DefaultTobinTax},
				{core.MicroMNTDenom, sdk.NewDecWithPrec(2, 2)}},
			SlashFraction:     slashFraction,
			SlashWindow:       slashWindow,
			MinValidPerWindow: minValidPerWindow,
		},
		[]types.ExchangeRatePrevote{},
		[]types.ExchangeRateVote{},
		map[string]sdk.Dec{},
		map[string]sdk.AccAddress{},
		map[string]int64{},
		[]types.AggregateExchangeRatePrevote{},
		[]types.AggregateExchangeRateVote{},
		map[string]sdk.Dec{},
	)

	fmt.Printf("Selected randomly generated oracle parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, oracleGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(oracleGenesis)
}
