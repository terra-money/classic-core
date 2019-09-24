package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Test a new validator entering the validator set
// Ensure that VotingInfo.StartHeight is set correctly
func TestHandleNewValidator(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	addr, val := ValAddrs[0], PubKeys[0]
	amt := sdk.TokensFromConsensusPower(100)
	sh := staking.NewHandler(input.StakingKeeper)

	// 1000 first blocks not a validator
	ctx := input.Ctx.WithBlockHeight(input.OracleKeeper.VotesWindow(input.Ctx) + 1)

	// Create a validator
	got := sh(ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.True(t, got.IsOK())
	staking.EndBlocker(ctx, input.StakingKeeper)

	require.Equal(
		t, input.BankKeeper.GetCoins(ctx, sdk.AccAddress(addr)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr).GetBondedTokens())

	// The validator miss one vote
	ballotAttendees := make(map[string]bool)
	ballotAttendees[addr.String()] = true
	input.OracleKeeper.HandleBallotSlashing(ctx, ballotAttendees)

	ctx = ctx.WithBlockHeight(input.OracleKeeper.VotesWindow(ctx) + 2)
	ballotAttendees[addr.String()] = false
	input.OracleKeeper.HandleBallotSlashing(ctx, ballotAttendees)

	info, found := input.OracleKeeper.getVotingInfo(ctx, addr)
	require.True(t, found)
	require.Equal(t, input.OracleKeeper.VotesWindow(ctx)+1, info.StartHeight)
	require.Equal(t, int64(2), info.IndexOffset)
	require.Equal(t, int64(1), info.MissedVotesCounter)

	// The validator should be bonded still, should not have been slashed
	validator := input.StakingKeeper.Validator(ctx, addr)
	require.Equal(t, sdk.Bonded, validator.GetStatus())
	bondPool := input.StakingKeeper.GetBondedPool(ctx)
	expTokens := sdk.TokensFromConsensusPower(100)
	require.Equal(t, expTokens.Int64(), bondPool.GetCoins().AmountOf(input.StakingKeeper.BondDenom(ctx)).Int64())
}

// Test slahsing a validator who did more than 5% wrong votes in VotingWindow
func TestSlash(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	addr, val := ValAddrs[0], PubKeys[0]
	amt := sdk.TokensFromConsensusPower(100)
	sk := input.StakingKeeper
	sh := staking.NewHandler(sk)
	ctx := input.Ctx

	got := sh(ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.True(t, got.IsOK())
	staking.EndBlocker(ctx, sk)

	height := int64(0)
	for ; height <= int64(50); height++ {
		ctx = ctx.WithBlockHeight(height)

		ballotAttendees := make(map[string]bool)
		ballotAttendees[addr.String()] = true
		input.OracleKeeper.HandleBallotSlashing(ctx, ballotAttendees)
	}

	// shouldn't be slashed
	validator := sk.Validator(ctx, addr)
	expTokens := sdk.TokensFromConsensusPower(100)
	require.Equal(t, expTokens.Int64(), validator.GetBondedTokens().Int64())

	// missed 95% blocks
	for ; height <= int64(1000); height++ {
		ctx = ctx.WithBlockHeight(height)

		ballotAttendees := make(map[string]bool)
		ballotAttendees[addr.String()] = false
		input.OracleKeeper.HandleBallotSlashing(ctx, ballotAttendees)
	}

	// shouldn't be slashed
	validator = sk.Validator(ctx, addr)
	require.Equal(t, expTokens.Int64(), validator.GetBondedTokens().Int64())

	for ; height <= int64(1001); height++ {
		ctx = ctx.WithBlockHeight(height)

		ballotAttendees := make(map[string]bool)
		ballotAttendees[addr.String()] = false
		input.OracleKeeper.HandleBallotSlashing(ctx, ballotAttendees)
	}

	// must be slashed
	validator = sk.Validator(ctx, addr)
	slashFraction := input.OracleKeeper.SlashFraction(ctx)
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(expTokens).TruncateInt(), validator.GetBondedTokens())
}
