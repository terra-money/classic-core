package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/types"
)

// Test a reward giving mechanism
func TestRewardBallotWinners(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr).GetBondedTokens())
	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr1)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr1).GetBondedTokens())

	// Add claim pools
	claim := types.NewClaim(10, 10, 0, addr)
	claim2 := types.NewClaim(20, 20, 0, addr1)
	claims := map[string]types.Claim{
		addr.String():  claim,
		addr1.String(): claim2,
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 30000000), sdk.NewInt64Coin(core.MicroUSDDenom, 40000000))
	acc := input.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	err = FundAccount(input, acc.GetAddress(), givingAmt)
	require.NoError(t, err)

	voteTargets := make(map[string]sdk.Dec)
	input.OracleKeeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) bool {
		voteTargets[denom] = tobinTax
		return false
	})

	votePeriodsPerWindow := sdk.NewDec((int64)(input.OracleKeeper.RewardDistributionWindow(input.Ctx))).
		QuoInt64((int64)(input.OracleKeeper.VotePeriod(input.Ctx))).
		TruncateInt64()
	input.OracleKeeper.RewardBallotWinners(ctx, (int64)(input.OracleKeeper.VotePeriod(input.Ctx)), (int64)(input.OracleKeeper.RewardDistributionWindow(input.Ctx)), voteTargets, claims)
	outstandingRewardsDec := input.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, addr)
	outstandingRewards, _ := outstandingRewardsDec.TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(core.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).TruncateInt(),
		outstandingRewards.AmountOf(core.MicroLunaDenom))
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(core.MicroUSDDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).TruncateInt(),
		outstandingRewards.AmountOf(core.MicroUSDDenom))

	outstandingRewardsDec1 := input.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, addr1)
	outstandingRewards1, _ := outstandingRewardsDec1.TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(core.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).MulInt64(2).TruncateInt(),
		outstandingRewards1.AmountOf(core.MicroLunaDenom))
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(core.MicroUSDDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).MulInt64(2).TruncateInt(),
		outstandingRewards1.AmountOf(core.MicroUSDDenom))
}
