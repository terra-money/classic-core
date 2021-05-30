package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/terra-money/core/x/market/internal/types"
)

// Simulation parameter constants
const (
	basePoolKey           = "base_pool"
	poolRecoveryPeriodKey = "pool_recovery_period"
	minStabilitySpreadKey = "min_spread"
)

// GenBasePool randomized BasePool
func GenBasePool(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(50000000000000).Add(sdk.NewDec(int64(r.Intn(10000000000))))
}

// GenPoolRecoveryPeriod randomized PoolRecoveryPeriod
func GenPoolRecoveryPeriod(r *rand.Rand) int64 {
	return int64(100 + r.Intn(10000000000))
}

// GenMinSpread randomized MinSpread
func GenMinSpread(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {

	var basePool sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, basePoolKey, &basePool, simState.Rand,
		func(r *rand.Rand) { basePool = GenBasePool(r) },
	)

	var poolRecoveryPeriod int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, poolRecoveryPeriodKey, &poolRecoveryPeriod, simState.Rand,
		func(r *rand.Rand) { poolRecoveryPeriod = GenPoolRecoveryPeriod(r) },
	)

	var minStabilitySpread sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, minStabilitySpreadKey, &minStabilitySpread, simState.Rand,
		func(r *rand.Rand) { minStabilitySpread = GenMinSpread(r) },
	)

	marketGenesis := types.NewGenesisState(
		sdk.ZeroDec(),
		types.Params{
			BasePool:           basePool,
			PoolRecoveryPeriod: poolRecoveryPeriod,
			MinStabilitySpread: minStabilitySpread,
		},
	)

	fmt.Printf("Selected randomly generated market parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, marketGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(marketGenesis)
}
