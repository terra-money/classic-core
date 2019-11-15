package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/oracle/internal/keeper"

	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestSlashAndResetMissCounters(t *testing.T) {
	input, _ := setup(t)

	slashWindow := input.OracleKeeper.SlashWindow(input.Ctx)
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidVotes := input.OracleKeeper.MinValidPerWindow(input.Ctx).MulInt64(slashWindow).TruncateInt64()
	// Case 1, no slash
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], slashWindow-minValidVotes)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// Case 2, slash
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], slashWindow-minValidVotes+1)
	SlashAndResetMissCounters(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt.Sub(slashFraction.MulInt(stakingAmt).TruncateInt()), validator.GetBondedTokens())
	require.True(t, validator.IsJailed())
}
