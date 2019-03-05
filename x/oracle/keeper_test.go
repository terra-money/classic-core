package oracle

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestPrice(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	// New context. There should be no price.
	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

	terraPrice := sdk.NewDecWithPrec(166, 2)
	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(terraPrice) && err == nil)
}

func TestVotes(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	votes := []PriceVote{}
	for i := 0; i < 3; i++ {
		vote := NewPriceVote(
			sdk.NewDecWithPrec(int64(i)+1, 0),
			assets.TerraDenom,
			sdk.OneInt(),
			addrs[i],
		)
		votes = append(votes, vote)
	}

	keeper.addVote(ctx, votes[0])

	// Should be one vote in total
	terraVotes := keeper.getVotes(ctx)[assets.TerraDenom]
	require.Equal(t, 1, len(votes))

	// Add the same vote; Should still be one vote in total
	keeper.addVote(ctx, votes[0])
	terraVotes = keeper.getVotes(ctx)[assets.TerraDenom]
	require.Equal(t, 1, len(votes))

	// Zero votes for an unrelated denom
	krwVotes := keeper.getVotes(ctx)[assets.KRWDenom]
	require.Equal(t, 0, len(krwVotes))

	// Should now be three votes in total
	keeper.addVote(ctx, votes[1])
	keeper.addVote(ctx, votes[2])
	terraVotes = keeper.getVotes(ctx)[assets.TerraDenom]
	require.Equal(t, 3, len(terraVotes))

	// Should now be two votes
	keeper.deleteVote(ctx, votes[0])
	terraVotes = keeper.getVotes(ctx)[assets.TerraDenom]
	require.Equal(t, 2, len(votes))

	// Should now be zero votes
	deleter := func(vote PriceVote) (stop bool) {
		keeper.deleteVote(ctx, vote)
		return false
	}
	keeper.iterateVotes(ctx, deleter)
	terraVotes = keeper.getVotes(ctx)[assets.TerraDenom]
	require.Equal(t, 0, len(votes))
}
