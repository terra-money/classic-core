package treasury

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/x/treasury/internal/keeper"
	"github.com/terra-money/core/x/treasury/internal/types"
)

func testTaxRateUpdateProposal(taxRate sdk.Dec) types.TaxRateUpdateProposal {
	return types.NewTaxRateUpdateProposal(
		"Test",
		"description",
		taxRate,
	)
}

func testRewardWeightUpdateProposal(rewardWeight sdk.Dec) types.RewardWeightUpdateProposal {
	return types.NewRewardWeightUpdateProposal(
		"Test",
		"description",
		rewardWeight,
	)
}

func TestTaxRateUpdateProposalHandler(t *testing.T) {
	input := keeper.CreateTestInput(t)

	taxRate := sdk.NewDecWithPrec(123, 5)
	tp := testTaxRateUpdateProposal(taxRate)
	hdlr := NewTreasuryPolicyUpdateHandler(input.TreasuryKeeper)
	require.NoError(t, hdlr(input.Ctx, tp))
	require.Equal(t, taxRate, input.TreasuryKeeper.GetTaxRate(input.Ctx))
}

func TestRewardWeightUpdateProposalHandler(t *testing.T) {
	input := keeper.CreateTestInput(t)

	rewardWeight := sdk.NewDecWithPrec(55, 3)
	tp := testRewardWeightUpdateProposal(rewardWeight)
	hdlr := NewTreasuryPolicyUpdateHandler(input.TreasuryKeeper)
	require.NoError(t, hdlr(input.Ctx, tp))
	require.Equal(t, rewardWeight, input.TreasuryKeeper.GetRewardWeight(input.Ctx))
}
