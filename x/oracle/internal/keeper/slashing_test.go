package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestSlashAndResetMissCounters(t *testing.T) {
	input := CreateTestInput(t)

	// Validator created
	amt := sdk.TokensFromConsensusPower(100)
	sh := staking.NewHandler(input.StakingKeeper)
	got := sh(input.Ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.True(t, got.IsOK())

	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	slashWindow := input.OracleKeeper.SlashWindow(input.Ctx)
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidVotes := input.OracleKeeper.MinValidPerWindow(input.Ctx).MulInt64(slashWindow).TruncateInt64()
	// Case 1, no slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], slashWindow-minValidVotes)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	validator := input.StakingKeeper.Validator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt, validator.GetBondedTokens())

	// Case 2, slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], slashWindow-minValidVotes+1)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator = input.StakingKeeper.Validator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt.Sub(slashFraction.MulInt(amt).TruncateInt()), validator.GetBondedTokens())
	require.True(t, validator.IsJailed())
}
