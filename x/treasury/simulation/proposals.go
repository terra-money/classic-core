package simulation

import (
	"math/rand"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-project/core/x/treasury/keeper"
	"github.com/terra-project/core/x/treasury/types"
)

// Simulation operation weights constants
const (
	OpWeightSubmitTaxRateUpdateProposal      = "op_weight_submit_tax_rate_update_proposal"
	OpWeightSubmitRewardWeightUpdateProposal = "op_weight_submit_reward_weight_update_proposal"
)

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSubmitTaxRateUpdateProposal,
			simappparams.DefaultWeightCommunitySpendProposal,
			SimulateTaxRateUpdateProposalContent(k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSubmitRewardWeightUpdateProposal,
			simappparams.DefaultWeightCommunitySpendProposal,
			SimulateRewardWeightUpdateProposalContent(k),
		),
	}
}

// SimulateTaxRateUpdateProposalContent generates random tax-rate-update proposal content
// nolint: funlen
func SimulateTaxRateUpdateProposalContent(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, _ []simtypes.Account) simtypes.Content {
		taxPolicy := k.TaxPolicy(ctx)
		diff := taxPolicy.RateMax.Sub(taxPolicy.RateMin)
		taxRate := simtypes.RandomDecAmount(r, diff).Add(taxPolicy.RateMin)

		return types.NewTaxRateUpdateProposal(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 100),
			taxRate,
		)
	}
}

// SimulateRewardWeightUpdateProposalContent generates random reward-weight-update proposal content
// nolint: funlen
func SimulateRewardWeightUpdateProposalContent(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, _ []simtypes.Account) simtypes.Content {
		rewardPolicy := k.RewardPolicy(ctx)
		diff := rewardPolicy.RateMax.Sub(rewardPolicy.RateMin)
		rewardWeight := simtypes.RandomDecAmount(r, diff).Add(rewardPolicy.RateMin)

		return types.NewTaxRateUpdateProposal(
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 100),
			rewardWeight,
		)
	}
}
