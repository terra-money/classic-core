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
	tp := keeper.GetPriceTarget(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(sdk.ZeroDec()))

	op := keeper.GetPriceObserved(ctx, assets.TerraDenom)
	require.True(t, op.Equal(sdk.ZeroDec()))

	terraTargetPrice := sdk.NewDecWithPrec(166, 2)
	keeper.setPriceTarget(ctx, assets.TerraDenom, terraTargetPrice)

	terraObservedPrice := sdk.NewDecWithPrec(174, 2)
	keeper.setPriceObserved(ctx, assets.TerraDenom, terraObservedPrice)

	tp = keeper.GetPriceTarget(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(terraTargetPrice))

	op = keeper.GetPriceObserved(ctx, assets.TerraDenom)
	require.True(t, op.Equal(terraObservedPrice))
}

func TestGetSetVotes(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	voteOne := PriceVote{
		FeedMsg: PriceFeedMsg{
			Denom:         assets.TerraDenom,
			TargetPrice:   sdk.OneDec(),
			ObservedPrice: sdk.OneDec(),
			Feeder:        addrs[0],
		},
		Power: sdk.NewDecWithPrec(10, 2),
	}

	voteTwo := PriceVote{
		FeedMsg: PriceFeedMsg{
			Denom:         assets.TerraDenom,
			TargetPrice:   sdk.OneDec(),
			ObservedPrice: sdk.OneDec(),
			Feeder:        addrs[1],
		},
		Power: sdk.NewDecWithPrec(10, 2),
	}

	voteThree := PriceVote{
		FeedMsg: PriceFeedMsg{
			Denom:         assets.TerraDenom,
			TargetPrice:   sdk.OneDec(),
			ObservedPrice: sdk.OneDec(),
			Feeder:        addrs[2],
		},
		Power: sdk.NewDecWithPrec(10, 2),
	}

	keeper.addVote(ctx, voteOne)

	// Should be one vote in total
	votes := keeper.getVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(votes))

	// Should still be one vote in total
	keeper.addVote(ctx, voteOne)
	votes = keeper.getVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(votes))

	// Zero votes for an unrelated denom
	keeper.addVote(ctx, voteOne)
	votes = keeper.getVotes(ctx, assets.KRWDenom)
	require.Equal(t, 0, len(votes))

	// Should now be three votes in total
	keeper.addVote(ctx, voteTwo)
	keeper.addVote(ctx, voteThree)
	votes = keeper.getVotes(ctx, assets.TerraDenom)
	require.Equal(t, 3, len(votes))

	// Should now be two votes
	keeper.deleteVote(ctx, voteOne)
	votes = keeper.getVotes(ctx, assets.TerraDenom)
	require.Equal(t, 2, len(votes))

	// Should now be zero votes
	deleter := func(vote PriceVote) (stop bool) {
		keeper.deleteVote(ctx, vote)
		return false
	}
	keeper.iterateVotes(ctx, assets.TerraDenom, deleter)
	votes = keeper.getVotes(ctx, assets.TerraDenom)
	require.Equal(t, 0, len(votes))
}
