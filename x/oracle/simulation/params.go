package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/oracle/internal/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyVotePeriod),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenVotePeriod(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyVoteThreshold),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenVoteThreshold(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyRewardBand),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenRewardBand(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyRewardDistributionWindow),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenRewardDistributionWindow(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeySlashFraction),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSlashFraction(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeySlashWindow),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenSlashWindow(r))
			},
		),
	}
}
