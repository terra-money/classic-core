package oracle

import (
	"encoding/hex"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/assets"

	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeperPrice(t *testing.T) {
	input := createTestInput(t)

	cnyPrice := sdk.NewDecWithPrec(839, int64(oracleDecPrecision)).MulInt64(assets.MicroUnit)
	gbpPrice := sdk.NewDecWithPrec(4995, int64(oracleDecPrecision)).MulInt64(assets.MicroUnit)
	krwPrice := sdk.NewDecWithPrec(2838, int64(oracleDecPrecision)).MulInt64(assets.MicroUnit)
	lunaPrice := sdk.NewDecWithPrec(3282384, int64(oracleDecPrecision)).MulInt64(assets.MicroUnit)

	// Set & get prices
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

func TestKeeperSwapPool(t *testing.T) {
	input := createTestInput(t)

	// Test AddSwapFeePool
	fees := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1000)))
	input.oracleKeeper.AddSwapFeePool(input.ctx, fees)

	// Test GetSwapFeePool
	feesQuery := input.oracleKeeper.GetSwapFeePool(input.ctx)
	require.Equal(t, fees, feesQuery)

	// Test clearSwapFeePool
	input.oracleKeeper.clearSwapFeePool(input.ctx)
	feesQuery = input.oracleKeeper.GetSwapFeePool(input.ctx)

	require.True(t, feesQuery.Empty())
}

func TestKeeperClaimPool(t *testing.T) {
	input := createTestInput(t)

	// Test addClaimPool
	claim := types.NewClaim(sdk.NewInt(10), addrs[0])
	claim2 := types.NewClaim(sdk.NewInt(20), addrs[1])
	claimPool := types.ClaimPool{claim, claim2}
	input.oracleKeeper.addClaimPool(input.ctx, claimPool)

	claim = types.NewClaim(sdk.NewInt(15), addrs[0])
	claim2 = types.NewClaim(sdk.NewInt(30), addrs[2])
	claimPool = types.ClaimPool{claim, claim2}
	input.oracleKeeper.addClaimPool(input.ctx, claimPool)

	// Test iterateClaimPool
	input.oracleKeeper.iterateClaimPool(input.ctx, func(recipient sdk.AccAddress, weight sdk.Int) (stop bool) {
		if recipient.Equals(addrs[0]) {
			require.Equal(t, sdk.NewInt(25), weight)
		} else if recipient.Equals(addrs[1]) {
			require.Equal(t, sdk.NewInt(20), weight)
		} else if recipient.Equals(addrs[2]) {
			require.Equal(t, sdk.NewInt(30), weight)
		}
		return false
	})
}
func TestKeeperVote(t *testing.T) {
	input := createTestInput(t)

	// Test addVote
	vote := NewPriceVote(sdk.OneDec(), assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	input.oracleKeeper.addVote(input.ctx, vote)

	// Test getVote
	voteQuery, err := input.oracleKeeper.getVote(input.ctx, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
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
	_, err = input.oracleKeeper.getVote(input.ctx, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.NotNil(t, err)
}

func TestKeeperPrevote(t *testing.T) {
	input := createTestInput(t)

	hash, _ := VoteHash("1234", sdk.OneDec(), assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	hexHas := hex.EncodeToString(hash)

	// Test addPrevote
	prevote := NewPricePrevote(hexHas, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]), 1)
	input.oracleKeeper.addPrevote(input.ctx, prevote)

	// Test getPrevote
	prevoteQuery, err := input.oracleKeeper.getPrevote(input.ctx, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)
	require.Equal(t, prevote, prevoteQuery)

	// Test iteratevotes
	input.oracleKeeper.iteratePrevotes(input.ctx, func(prevote PricePrevote) bool {
		require.Equal(t, prevote, prevoteQuery)
		return true
	})

	// Test deletevote
	input.oracleKeeper.deletePrevote(input.ctx, prevote)
	_, err = input.oracleKeeper.getPrevote(input.ctx, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.NotNil(t, err)
}

func TestKeeperParams(t *testing.T) {
	input := createTestInput(t)

	// Test default params setting
	input.oracleKeeper.SetParams(input.ctx, DefaultParams())
	params := input.oracleKeeper.GetParams(input.ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := int64(10)
	voteThreshold := sdk.NewDecWithPrec(1, 10)
	oracleRewardBand := sdk.NewDecWithPrec(1, 2)

	// Should really test validateParams, but skipping because obvious
	newParams := NewParams(votePeriod, voteThreshold, oracleRewardBand)
	input.oracleKeeper.SetParams(input.ctx, newParams)

	storedParams := input.oracleKeeper.GetParams(input.ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, newParams, storedParams)
}

func TestKeeperFeederDelegation(t *testing.T) {
	input := createTestInput(t)

	// Test default getters and setters
	delegate := input.oracleKeeper.GetFeedDelegate(input.ctx, sdk.ValAddress(addrs[0]))
	require.Equal(t, delegate, addrs[0])

	input.oracleKeeper.SetFeedDelegate(input.ctx, sdk.ValAddress(addrs[0]), addrs[1])
	delegate = input.oracleKeeper.GetFeedDelegate(input.ctx, sdk.ValAddress(addrs[0]))
	require.Equal(t, delegate, addrs[1])
}
