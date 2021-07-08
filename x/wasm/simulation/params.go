package simulation

//DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/wasm/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxContractSize),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractSize(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxContractGas),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractGas(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxContractMsgSize),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractMsgSize(r))
			},
		),
	}
}
