package simulation

import (
	"math/rand"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/treasury/internal/keeper"
	"github.com/terra-money/core/x/treasury/internal/types"
)

// Simulation operation weights constants
const (
	OpWeightSubmitTaxRateUpdateProposal      = "op_weight_submit_tax_rate_update_proposal"
	OpWeightSubmitRewardWeightUpdateProposal = "op_weight_submit_reward_weight_update_proposal"
)

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitTaxRateUpdateProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateTaxRateUpdateProposalContent(k),
		},
		{
			AppParamsKey:       OpWeightSubmitRewardWeightUpdateProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateRewardWeightUpdateProposalContent(k),
		},
	}
}

// SimulateTaxRateUpdateProposalContent generates random tax-rate-update proposal content
// nolint: funlen
func SimulateTaxRateUpdateProposalContent(k keeper.Keeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, _ []simulation.Account) govtypes.Content {
		taxPolicy := k.TaxPolicy(ctx)
		diff := taxPolicy.RateMax.Sub(taxPolicy.RateMin)
		taxRate := simulation.RandomDecAmount(r, diff).Add(taxPolicy.RateMin)

		return types.NewTaxRateUpdateProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			taxRate,
		)
	}
}

// SimulateRewardWeightUpdateProposalContent generates random reward-weight-update proposal content
// nolint: funlen
func SimulateRewardWeightUpdateProposalContent(k keeper.Keeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, _ []simulation.Account) govtypes.Content {
		rewardPolicy := k.RewardPolicy(ctx)
		diff := rewardPolicy.RateMax.Sub(rewardPolicy.RateMin)
		rewardWeight := simulation.RandomDecAmount(r, diff).Add(rewardPolicy.RateMin)

		return types.NewTaxRateUpdateProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			rewardWeight,
		)
	}
}
