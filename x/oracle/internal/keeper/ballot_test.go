package keeper

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestOrganize(t *testing.T) {
	input := CreateTestInput(t)

	sdrBallot := types.ExchangeRateBallot{
		types.NewExchangeRateVote(sdk.NewDec(17), core.MicroSDRDenom, ValAddrs[0]),
		types.NewExchangeRateVote(sdk.NewDec(10), core.MicroSDRDenom, ValAddrs[1]),
		types.NewExchangeRateVote(sdk.NewDec(6), core.MicroSDRDenom, ValAddrs[2]),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewExchangeRateVote(sdk.NewDec(1000), core.MicroKRWDenom, ValAddrs[0]),
		types.NewExchangeRateVote(sdk.NewDec(1300), core.MicroKRWDenom, ValAddrs[1]),
		types.NewExchangeRateVote(sdk.NewDec(2000), core.MicroKRWDenom, ValAddrs[2]),
	}

	for _, vote := range sdrBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote)
	}
	for _, vote := range krwBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote)
	}

	// oranize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[core.MicroSDRDenom])
	sort.Sort(ballotMap[core.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[core.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[core.MicroKRWDenom])

}
