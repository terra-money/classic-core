package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/oracle/internal/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestSlashAndResetMissCounters(t *testing.T) {
	input, _ := setup(t)

	votePeriodsPerWindow := input.OracleKeeper.SlashWindow(input.Ctx)
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidVotes := input.OracleKeeper.MinValidPerWindow(input.Ctx).MulInt64(votePeriodsPerWindow).TruncateInt64()
	// Case 1, no slash
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], votePeriodsPerWindow-minValidVotes)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	validator, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// Case 2, slash
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt.Sub(slashFraction.MulInt(stakingAmt).TruncateInt()), validator.GetBondedTokens())
	require.True(t, validator.IsJailed())

	// Case 3, slash unbonded validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	validator.Status = sdk.Unbonded
	validator.Jailed = false
	validator.Tokens = stakingAmt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.Tokens)
	require.False(t, validator.IsJailed())

	// Case 4, slash jailed validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	validator.Status = sdk.Bonded
	validator.Jailed = true
	validator.Tokens = stakingAmt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.Tokens)
}
