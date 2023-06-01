package simulation

import (
	"fmt"
	"math/rand" //#nosec G404

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/classic-terra/core/v2/x/feeshare/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyEnableFeeShare),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%v\"", GenEnableFeeShare(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyDeveloperShares),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenDeveloperShares(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyAllowedDenoms),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenAllowedDenoms(r))
			},
		),
	}
}
