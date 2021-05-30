package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"

	"github.com/terra-money/core/x/treasury/internal/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyTaxPolicy),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenTaxPolicy(r))
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyRewardPolicy),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenRewardPolicy(r))
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeySeigniorageBurdenTarget),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSeigniorageBurdenTarget(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMiningIncrement),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMiningIncrement(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyWindowShort),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowShort(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyWindowLong),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowLong(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyWindowProbation),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowProbation(r))
			},
		),
	}
}
