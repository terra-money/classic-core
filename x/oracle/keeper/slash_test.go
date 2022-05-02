package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/terra-money/core/x/oracle/types"
)

func TestSlashAndResetMissCounters(t *testing.T) {
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

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidVotes := input.OracleKeeper.MinValidPerWindow(input.Ctx).MulInt64(votePeriodsPerWindow).TruncateInt64()

	// clear event
	input.Ctx = input.Ctx.WithEventManager(sdk.NewEventManager())

	// Case 1, no slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], uint64(votePeriodsPerWindow-minValidVotes))
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	validator, _ := input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt, validator.GetBondedTokens())

	// Event should be empty
	require.Zero(t, len(input.Ctx.EventManager().Events()))

	// Preload power and address before slashing
	consAddr, err := validator.GetConsAddr()
	require.NoError(t, err)

	powerReduction := input.StakingKeeper.PowerReduction(ctx)
	power := validator.GetConsensusPower(powerReduction)

	// Case 2, slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], uint64(votePeriodsPerWindow-minValidVotes+1))
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt.Sub(slashFraction.MulInt(amt).TruncateInt()), validator.GetBondedTokens())
	require.True(t, validator.IsJailed())

	// Event should be exists
	require.Equal(t,
		sdk.NewEvent(
			slashingtypes.EventTypeSlash,
			sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(slashingtypes.AttributeKeyReason, types.AttributeValueMissingOracleVote),
			sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consAddr.String()),
		),
		input.Ctx.EventManager().Events()[0],
	)

	// Case 3, slash unbonded validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	validator.Status = stakingtypes.Unbonded
	validator.Jailed = false
	validator.Tokens = amt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], uint64(votePeriodsPerWindow-minValidVotes+1))
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt, validator.Tokens)
	require.False(t, validator.IsJailed())

	// Case 4, slash jailed validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	validator.Status = stakingtypes.Bonded
	validator.Jailed = true
	validator.Tokens = amt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], uint64(votePeriodsPerWindow-minValidVotes+1))
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, amt, validator.Tokens)
}
