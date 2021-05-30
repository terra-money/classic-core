package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/market/internal/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyBasePool),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBasePool(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyPoolRecoveryPeriod),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenPoolRecoveryPeriod(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMinStabilitySpread),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMinSpread(r))
			},
		),
	}
}
