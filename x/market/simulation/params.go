package simulation

//DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/market/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMintBasePool),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMintBasePool(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyBurnBasePool),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBurnBasePool(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPoolRecoveryPeriod),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenPoolRecoveryPeriod(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMinStabilitySpread),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMinSpread(r))
			},
		),
	}
}
