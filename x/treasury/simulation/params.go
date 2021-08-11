package simulation

//DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/treasury/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyTaxPolicy),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenTaxPolicy(r))
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyRewardPolicy),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenRewardPolicy(r))
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySeigniorageBurdenTarget),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSeigniorageBurdenTarget(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMiningIncrement),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMiningIncrement(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyWindowShort),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowShort(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyWindowLong),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowLong(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyWindowProbation),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenWindowProbation(r))
			},
		),
	}
}
