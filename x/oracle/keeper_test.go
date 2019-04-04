package oracle

import (
	"github.com/terra-project/core/types/assets"

	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeperPrice(t *testing.T) {
	input := createTestInput(t)

	cnyPrice := sdk.NewDecWithPrec(839, precision).MulInt64(assets.MicroUnit)
	gbpPrice := sdk.NewDecWithPrec(4995, precision).MulInt64(assets.MicroUnit)
	krwPrice := sdk.NewDecWithPrec(2838, precision).MulInt64(assets.MicroUnit)
	lunaPrice := sdk.NewDecWithPrec(3282384, precision).MulInt64(assets.MicroUnit)

	// Set prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, cnyPrice)
	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroCNYDenom)
	require.Nil(t, err)
	require.Equal(t, cnyPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroGBPDenom, gbpPrice)
	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroGBPDenom)
	require.Nil(t, err)
	require.Equal(t, gbpPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, krwPrice)
	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroKRWDenom)
	require.Nil(t, err)
	require.Equal(t, krwPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroLunaDenom, lunaPrice)
	price, _ = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroLunaDenom)
	require.Equal(t, sdk.OneDec(), price)
}

func TestKeeperVote(t *testing.T) {
	input := createTestInput(t)

	// Test addvote
	vote := NewPriceVote(sdk.OneDec(), assets.MicroSDRDenom, sdk.NewInt(3458).MulRaw(assets.MicroUnit), addrs[0])
	input.oracleKeeper.addVote(input.ctx, vote)

	// Test getVote
	voteQuery, err := input.oracleKeeper.getVote(input.ctx, assets.MicroSDRDenom, addrs[0])
	require.Nil(t, err)
	require.Equal(t, vote, voteQuery)

	// Test iteratevotes
	input.oracleKeeper.iterateVotes(input.ctx, func(vote PriceVote) bool {
		require.Equal(t, vote, voteQuery)
		return true
	})

	// Test collectvotes
	votes := input.oracleKeeper.collectVotes(input.ctx)
	require.True(t, len(votes) == 1)
	require.True(t, len(votes[assets.MicroSDRDenom]) == 1)
	require.Equal(t, vote, votes[assets.MicroSDRDenom][0])

	// Test deletevote
	input.oracleKeeper.deleteVote(input.ctx, vote)
	_, err = input.oracleKeeper.getVote(input.ctx, assets.MicroSDRDenom, addrs[0])
	require.NotNil(t, err)
}

func TestKeeperDropCounter(t *testing.T) {
	input := createTestInput(t)

	for i := 1; i < 40; i++ {
		counter := input.oracleKeeper.incrementDropCounter(input.ctx, assets.MicroSDRDenom)
		require.Equal(t, sdk.NewInt(int64(i)), counter)
	}

	input.oracleKeeper.resetDropCounter(input.ctx, assets.MicroSDRDenom)
	store := input.ctx.KVStore(input.oracleKeeper.key)
	b := store.Get(keyDropCounter(assets.MicroSDRDenom))
	require.Nil(t, b)
}

func TestKeeperParams(t *testing.T) {
	input := createTestInput(t)

	// Test default params setting
	input.oracleKeeper.SetParams(input.ctx, DefaultParams())
	params := input.oracleKeeper.GetParams(input.ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := sdk.NewInt(10)
	voteThreshold := sdk.NewDecWithPrec(1, 10)
	dropThreshold := sdk.NewInt(10)

	// Should really test validateParams, but skipping because obvious
	newParams := NewParams(votePeriod, voteThreshold, dropThreshold)
	input.oracleKeeper.SetParams(input.ctx, newParams)

	storedParams := input.oracleKeeper.GetParams(input.ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, newParams, storedParams)
}
