package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/wasm/internal/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMaxContractSize),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractSize(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMaxContractGas),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractGas(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMaxContractMsgSize),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxContractMsgSize(r))
			},
		),
	}
}
