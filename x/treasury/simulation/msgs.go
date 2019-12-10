package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-project/core/x/treasury"
	"github.com/terra-project/core/x/treasury/internal/types"
)

// SimulateTaxRateUpdateProposalContent generates random tax-rate-update proposal content
func SimulateTaxRateUpdateProposalContent(k treasury.Keeper) govsim.ContentSimulator {
	return func(r *rand.Rand, _ *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) gov.Content {

		targetTaxRate := sdk.NewDecWithPrec(r.Int63n(100), 2)
		if targetTaxRate.GT(types.DefaultTaxPolicy.RateMax) {
			targetTaxRate = types.DefaultTaxPolicy.RateMax
		}

		if targetTaxRate.LT(types.DefaultTaxPolicy.RateMin) {
			targetTaxRate = types.DefaultTaxPolicy.RateMin
		}

		return treasury.NewTaxRateUpdateProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			targetTaxRate,
		)
	}
}

// SimulateRewardWeightUpdateProposalContent generates random tax-rate-update proposal content
func SimulateRewardWeightUpdateProposalContent(k treasury.Keeper) govsim.ContentSimulator {
	return func(r *rand.Rand, _ *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) gov.Content {

		targetRewardWeight := sdk.NewDecWithPrec(r.Int63n(100), 2)
		if targetRewardWeight.GT(types.DefaultRewardPolicy.RateMax) {
			targetRewardWeight = types.DefaultRewardPolicy.RateMax
		}

		if targetRewardWeight.LT(types.DefaultRewardPolicy.RateMin) {
			targetRewardWeight = types.DefaultRewardPolicy.RateMin
		}

		return treasury.NewRewardWeightUpdateProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			targetRewardWeight,
		)
	}
}
