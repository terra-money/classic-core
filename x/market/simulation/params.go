package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
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
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyMinSpread),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMinSpread(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParmaStoreKeyTobinTax),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenTobinTax(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParmaStoreKeyIlliquidTobinTaxList),
			func(r *rand.Rand) string {
				return fmt.Sprintf("[{\"denom\": \"%s\", \"tax_rate\": \"%s\"}]", core.MicroMNTDenom, GenIlliquidTobinTaxRate(r))
			},
		),
	}
}
