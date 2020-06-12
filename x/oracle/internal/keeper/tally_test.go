package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestTallyCrossRate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	valiCount := 5
	for i := 0; i < valiCount; i++ {
		_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[i], PubKeys[i], amt))
		require.NoError(t, err)
	}
	staking.EndBlocker(ctx, input.StakingKeeper)

	// Set WhiteList, TobinTax
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{
		{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax},
		{Name: core.MicroSDRDenom, TobinTax: types.DefaultTobinTax},
		{Name: core.MicroUSDDenom, TobinTax: types.DefaultTobinTax},
		{Name: core.MicroMNTDenom, TobinTax: types.DefaultTobinTax},
	}
	input.OracleKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.ClearTobinTaxes(ctx)
	for _, denom := range params.Whitelist {
		input.OracleKeeper.SetTobinTax(ctx, denom.Name, denom.TobinTax)
	}

	usdBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(100, int64(OracleDecPrecision)), core.MicroUSDDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(100, int64(OracleDecPrecision)), core.MicroUSDDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(100, int64(OracleDecPrecision)), core.MicroUSDDenom, ValAddrs[2]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(100, int64(OracleDecPrecision)), core.MicroUSDDenom, ValAddrs[3]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(20000, int64(OracleDecPrecision)), core.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(25000, int64(OracleDecPrecision)), core.MicroKRWDenom, ValAddrs[2]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(30000, int64(OracleDecPrecision)), core.MicroKRWDenom, ValAddrs[3]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(35000, int64(OracleDecPrecision)), core.MicroKRWDenom, ValAddrs[4]), power),
	}

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(75, int64(OracleDecPrecision)), core.MicroSDRDenom, ValAddrs[2]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(80, int64(OracleDecPrecision)), core.MicroSDRDenom, ValAddrs[3]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(75, int64(OracleDecPrecision)), core.MicroSDRDenom, ValAddrs[4]), power),
	}

	mntBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(20000, int64(OracleDecPrecision)), core.MicroMNTDenom, ValAddrs[2]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(40000, int64(OracleDecPrecision)), core.MicroMNTDenom, ValAddrs[3]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(40000, int64(OracleDecPrecision)), core.MicroMNTDenom, ValAddrs[4]), power),
	}

	// Ignored, not in DefaultWhitelist
	eurBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(85, int64(OracleDecPrecision)), core.MicroEURDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(85, int64(OracleDecPrecision)), core.MicroEURDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(85, int64(OracleDecPrecision)), core.MicroEURDenom, ValAddrs[2]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(85, int64(OracleDecPrecision)), core.MicroEURDenom, ValAddrs[3]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDecWithPrec(85, int64(OracleDecPrecision)), core.MicroEURDenom, ValAddrs[4]), power),
	}


	for _, vote := range usdBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range krwBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range sdrBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range mntBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range eurBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	voteTargets := make(map[string]sdk.Dec)
	input.OracleKeeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) bool {
		voteTargets[denom] = tobinTax
		return false
	})

	crossRates := input.OracleKeeper.TallyCrossRate(input.Ctx, ballotMap, voteTargets)
	require.Equal(t, 4, len(crossRates))
	for _, cer := range crossRates {
		input.OracleKeeper.SetCrossExchangeRate(ctx, cer)
	}
	cer1, err := input.OracleKeeper.GetCrossExchangeRate(input.Ctx, core.MicroKRWDenom, core.MicroMNTDenom)
	require.NoError(t, err)
	require.Equal(t, core.MicroKRWDenom, cer1.Denom1)
	require.Equal(t, core.MicroMNTDenom, cer1.Denom2)
	require.Equal(t, sdk.NewDecWithPrec(875, 3), cer1.CrossExchangeRate)

	cer2, err := input.OracleKeeper.GetCrossExchangeRate(input.Ctx, core.MicroKRWDenom, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, core.MicroKRWDenom, cer2.Denom1)
	require.Equal(t, core.MicroSDRDenom, cer2.Denom2)
	require.Equal(t, sdk.NewDec(375), cer2.CrossExchangeRate)

	cer3, err := input.OracleKeeper.GetCrossExchangeRate(input.Ctx, core.MicroKRWDenom, core.MicroUSDDenom)
	require.NoError(t, err)
	require.Equal(t, core.MicroKRWDenom, cer3.Denom1)
	require.Equal(t, core.MicroUSDDenom, cer3.Denom2)
	require.Equal(t, sdk.NewDec(250), cer3.CrossExchangeRate)

	cer4, err := input.OracleKeeper.GetCrossExchangeRate(input.Ctx, core.MicroMNTDenom, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, core.MicroMNTDenom, cer4.Denom1)
	require.Equal(t, core.MicroSDRDenom, cer4.Denom2)
	require.Equal(t, sdk.NewDec(500), cer4.CrossExchangeRate)

}
