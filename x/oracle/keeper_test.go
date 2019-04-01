package oracle

import (
	"terra/types/assets"

	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeperPrice(t *testing.T) {
	input := createTestInput()

	cnyPrice := sdk.NewDecWithPrec(839, precision)
	gbpPrice := sdk.NewDecWithPrec(4995, precision)
	krwPrice := sdk.NewDecWithPrec(2838, precision)
	lunaPrice := sdk.NewDecWithPrec(3282384, precision)

	// Set prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.CNYDenom, cnyPrice)
	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.CNYDenom)
	require.Nil(t, err)
	require.Equal(t, cnyPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.GBPDenom, gbpPrice)
	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.GBPDenom)
	require.Nil(t, err)
	require.Equal(t, gbpPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.KRWDenom, krwPrice)
	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.KRWDenom)
	require.Nil(t, err)
	require.Equal(t, krwPrice, price)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.LunaDenom, lunaPrice)
	price, _ = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.LunaDenom)
	require.Equal(t, sdk.OneDec(), price)
}

func TestKeeperVote(t *testing.T) {
	input := createTestInput()

	// Test addvote
	vote := NewPriceVote(sdk.OneDec(), assets.SDRDenom, sdk.NewInt(3458), addrs[0])
	input.oracleKeeper.addVote(input.ctx, vote)

	// Test getVote
	voteQuery, err := input.oracleKeeper.getVote(input.ctx, assets.SDRDenom, addrs[0])
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
	require.True(t, len(votes[assets.SDRDenom]) == 1)
	require.Equal(t, vote, votes[assets.SDRDenom][0])

	// Test deletevote
	input.oracleKeeper.deleteVote(input.ctx, vote)
	_, err = input.oracleKeeper.getVote(input.ctx, assets.SDRDenom, addrs[0])
	require.NotNil(t, err)
}

func TestKeeperDropCounter(t *testing.T) {
	input := createTestInput()

	for i := 1; i < 40; i++ {
		counter := input.oracleKeeper.incrementDropCounter(input.ctx, assets.SDRDenom)
		require.Equal(t, sdk.NewInt(int64(i)), counter)
	}

	input.oracleKeeper.resetDropCounter(input.ctx, assets.SDRDenom)
	store := input.ctx.KVStore(input.oracleKeeper.key)
	b := store.Get(keyDropCounter(assets.SDRDenom))
	require.Nil(t, b)
}

func TestKeeperParams(t *testing.T) {
	input := createTestInput()

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
