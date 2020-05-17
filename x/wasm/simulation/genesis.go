package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// Simulation parameter constants
const (
	maxContractSizeKey = "max_contract_size"
	maxContractGasKey  = "max_contract_gas"
	gasMultiplierKey   = "gas_multiplier"
)

// GenMaxContractSize randomized MaxContractSize
func GenMaxContractSize(r *rand.Rand) int64 {
	return int64(1024 + r.Intn(499*1024))
}

// GenMaxContractGas randomized MaxContractGas
func GenMaxContractGas(r *rand.Rand) uint64 {
	return uint64(10000 + r.Intn(500_000_000))
}

// GenGasMultiplier randomized GasMultiplier
func GenGasMultiplier(r *rand.Rand) uint64 {
	return uint64(1 + r.Intn(99))
}

// RandomizedGenState generates a random GenesisState for wasm
func RandomizedGenState(simState *module.SimulationState) {

	var maxContractSize int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractSizeKey, &maxContractSize, simState.Rand,
		func(r *rand.Rand) { maxContractSize = GenMaxContractSize(r) },
	)

	var maxContractGas uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractGasKey, &maxContractGas, simState.Rand,
		func(r *rand.Rand) { maxContractGas = GenMaxContractGas(r) },
	)

	var gasMultiplier uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, gasMultiplierKey, &gasMultiplier, simState.Rand,
		func(r *rand.Rand) { gasMultiplier = GenGasMultiplier(r) },
	)

	wasmGenesis := types.NewGenesisState(
		types.Params{
			MaxContractSize: maxContractSize,
			MaxContractGas:  maxContractGas,
			GasMultiplier:   gasMultiplier,
		},
		0,
		0,
		[]types.Code{},
		[]types.Contract{},
	)

	fmt.Printf("Selected randomly generated wasm parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, wasmGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(wasmGenesis)
}
