package simulation

//DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/terra-money/core/x/wasm/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	maxContractSizeKey     = "max_contract_size"
	maxContractGasKey      = "max_contract_gas"
	maxContractMsgSizeKey  = "max_contract_msg_size"
	maxContractDataSizeKey = "max_contract_data_size"
	EventParamsKey         = "event_params"
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

// GenMaxContractDataSize randomized MaxContractDataSize
func GenMaxContractDataSize(r *rand.Rand) uint64 {
	return uint64(256 + r.Intn(512))
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

	var maxContractDataSize uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxContractDataSizeKey, &maxContractDataSize, simState.Rand,
		func(r *rand.Rand) { maxContractDataSize = GenMaxContractDataSize(r) },
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

	bz, err := json.MarshalIndent(&wasmGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected randomly generated wasm parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(wasmGenesis)
}
