package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/terra-money/core/x/wasm/internal/types"
)

// Simulation parameter constants
const (
	maxContractSizeKey    = "max_contract_size"
	maxContractGasKey     = "max_contract_gas"
	maxContractMsgSizeKey = "max_contract_msg_size"
	gasMultiplierKey      = "gas_multiplier"
)

// GenMaxContractSize randomized MaxContractSize
func GenMaxContractSize(r *rand.Rand) uint64 {
	return uint64(300*1024 + r.Intn(200*1024))
}

// GenMaxContractGas randomized MaxContractGas
func GenMaxContractGas(r *rand.Rand) uint64 {
	return uint64(10_000_000 + r.Intn(90_000_000))
}

// GenMaxContractMsgSize randomized MaxContractMsgSize
func GenMaxContractMsgSize(r *rand.Rand) uint64 {
	return uint64(128 + r.Intn(9*1024))
}

// GenGasMultiplier randomized GasMultiplier
func GenGasMultiplier(r *rand.Rand) uint64 {
	return uint64(1 + r.Intn(99))
}

// RandomizedGenState generates a random GenesisState for wasm
func RandomizedGenState(simState *module.SimulationState) {

	var maxContractSize uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractSizeKey, &maxContractSize, simState.Rand,
		func(r *rand.Rand) { maxContractSize = GenMaxContractSize(r) },
	)

	var maxContractGas uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractGasKey, &maxContractGas, simState.Rand,
		func(r *rand.Rand) { maxContractGas = GenMaxContractGas(r) },
	)

	var maxContractMsgSize uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractMsgSizeKey, &maxContractMsgSize, simState.Rand,
		func(r *rand.Rand) { maxContractMsgSize = GenMaxContractMsgSize(r) },
	)

	var gasMultiplier uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, gasMultiplierKey, &gasMultiplier, simState.Rand,
		func(r *rand.Rand) { gasMultiplier = GenGasMultiplier(r) },
	)

	wasmGenesis := types.NewGenesisState(
		types.Params{
			MaxContractSize:    maxContractSize,
			MaxContractGas:     maxContractGas,
			MaxContractMsgSize: maxContractMsgSize,
		},
		0,
		0,
		[]types.Code{},
		[]types.Contract{},
	)

	fmt.Printf("Selected randomly generated wasm parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, wasmGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(wasmGenesis)
}
