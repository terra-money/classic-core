package keeper

import (
	"math/rand"
	"testing"

	core "github.com/terra-project/core/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSettle(t *testing.T) {
	input := CreateTestInput(t)

	burnAmt := sdk.NewInt(rand.Int63() + 1)
	supply := input.BankKeeper.GetSupply(input.Ctx)
	supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, burnAmt)))
	input.BankKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(int64(core.BlocksPerWeek))
	supply.SetTotal(sdk.NewCoins())
	input.BankKeeper.SetSupply(input.Ctx, supply)

	// check seigniorage update
	require.Equal(t, burnAmt, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx))

	input.TreasuryKeeper.SettleSeigniorage(input.Ctx)
	supply = input.BankKeeper.GetSupply(input.Ctx)
	feePool := input.DistrKeeper.GetFeePool(input.Ctx)

	// Reward weight portion of seigniorage burned
	rewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	communityPoolAmt := burnAmt.Sub(rewardWeight.MulInt(burnAmt).TruncateInt())

	require.Equal(t, communityPoolAmt, supply.GetTotal().AmountOf(core.MicroLunaDenom))
	require.Equal(t, communityPoolAmt, feePool.CommunityPool.AmountOf(core.MicroLunaDenom).TruncateInt())
}
