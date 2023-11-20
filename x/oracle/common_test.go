package oracle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/classic-terra/core/v2/x/oracle"
	"github.com/classic-terra/core/v2/x/oracle/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	stakingAmt = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	randomExchangeRate = sdk.NewDec(1700)
)

func setupWithSmallVotingPower(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := oracle.NewHandler(input.OracleKeeper)

	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(input.StakingKeeper)
	_, err := stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.ValPubKeys[0], sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)))
	require.NoError(t, err)

	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := oracle.NewHandler(input.OracleKeeper)

	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(input.StakingKeeper)

	// Validator created
	_, err := stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.ValPubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.ValPubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.ValPubKeys[2], stakingAmt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}

func setupVal5(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := oracle.NewHandler(input.OracleKeeper)

	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(input.StakingKeeper)

	// Validator created
	_, err := stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.ValPubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.ValPubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.ValPubKeys[2], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[3], keeper.ValPubKeys[3], stakingAmt))
	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[4], keeper.ValPubKeys[4], stakingAmt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}
