package simulation

//DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/terra-money/core/x/market/types"
)

// Simulation parameter constants
const (
	mintBasePoolKey       = "mint_base_pool"
	burnBasePoolKey       = "burn_base_pool"
	poolRecoveryPeriodKey = "pool_recovery_period"
	minStabilitySpreadKey = "min_spread"
)

// GenMintBasePool randomized MintBasePool
func GenMintBasePool(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(50000000000000).Add(sdk.NewDec(int64(r.Intn(10000000000))))
}

// GenBurnBasePool randomized BurnBasePool
func GenBurnBasePool(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(50000000000000).Add(sdk.NewDec(int64(r.Intn(10000000000))))
}

// GenPoolRecoveryPeriod randomized PoolRecoveryPeriod
func GenPoolRecoveryPeriod(r *rand.Rand) uint64 {
	return uint64(100 + r.Intn(10000000000))
}

// GenMinSpread randomized MinSpread
func GenMinSpread(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(100)), 3))
}

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {

	var mintBasePool sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintBasePoolKey, &mintBasePool, simState.Rand,
		func(r *rand.Rand) { mintBasePool = GenMintBasePool(r) },
	)

	var burnBasePool sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, burnBasePoolKey, &burnBasePool, simState.Rand,
		func(r *rand.Rand) { burnBasePool = GenBurnBasePool(r) },
	)

	var poolRecoveryPeriod uint64
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
		sdk.ZeroDec(),
		types.Params{
			MintBasePool:       mintBasePool,
			BurnBasePool:       burnBasePool,
			PoolRecoveryPeriod: poolRecoveryPeriod,
			MinStabilitySpread: minStabilitySpread,
		},
	)

	bz, err := json.MarshalIndent(&marketGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected randomly generated market parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(marketGenesis)
}
