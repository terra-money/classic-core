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

func TestTargetVotes(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	votes := []PriceVote{}
	for i := 0; i < 3; i++ {
		vote := NewPriceVote(
			sdk.NewDecWithPrec(int64(i)+1, 0),
			assets.TerraDenom,
			sdk.NewDecWithPrec(10, 2),
			addrs[i],
		)
		votes = append(votes, vote)
	}

	keeper.addTargetVote(ctx, votes[0])

	// Should be one TargetVote in total
	targetVotes := keeper.getTargetVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(targetVotes))

	// Should still be one TargetVote in total
	keeper.addTargetVote(ctx, votes[0])
	targetVotes = keeper.getTargetVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(targetVotes))

	// Zero TargetVotes for an unrelated denom
	keeper.addTargetVote(ctx, votes[0])
	targetVotes = keeper.getTargetVotes(ctx, assets.KRWDenom)
	require.Equal(t, 0, len(targetVotes))

	// Should now be three TargetVotes in total
	keeper.addTargetVote(ctx, votes[1])
	keeper.addTargetVote(ctx, votes[2])
	targetVotes = keeper.getTargetVotes(ctx, assets.TerraDenom)
	require.Equal(t, 3, len(targetVotes))

	// Should now be two TargetVotes
	keeper.deleteTargetVote(ctx, votes[0])
	targetVotes = keeper.getTargetVotes(ctx, assets.TerraDenom)
	require.Equal(t, 2, len(targetVotes))

	// Should now be zero TargetVotes
	deleter := func(TargetVote PriceVote) (stop bool) {
		keeper.deleteTargetVote(ctx, TargetVote)
		return false
	}
	keeper.iterateTargetVotes(ctx, assets.TerraDenom, deleter)
	targetVotes = keeper.getTargetVotes(ctx, assets.TerraDenom)
	require.Equal(t, 0, len(targetVotes))
}

func TestObservedVotes(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 3)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	votes := []PriceVote{}
	for i := 0; i < 3; i++ {
		vote := NewPriceVote(
			sdk.NewDecWithPrec(int64(i)+1, 0),
			assets.TerraDenom,
			sdk.NewDecWithPrec(10, 2),
			addrs[i],
		)
		votes = append(votes, vote)
	}

	keeper.addObservedVote(ctx, votes[0])

	// Should be one ObservedVote in total
	observedVotes := keeper.getObservedVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(observedVotes))

	// Should still be one ObservedVote in total
	keeper.addObservedVote(ctx, votes[0])
	observedVotes = keeper.getObservedVotes(ctx, assets.TerraDenom)
	require.Equal(t, 1, len(observedVotes))

	// Zero ObservedVotes for an unrelated denom
	keeper.addObservedVote(ctx, votes[0])
	observedVotes = keeper.getObservedVotes(ctx, assets.KRWDenom)
	require.Equal(t, 0, len(observedVotes))

	// Should now be three ObservedVotes in total
	keeper.addObservedVote(ctx, votes[1])
	keeper.addObservedVote(ctx, votes[2])
	observedVotes = keeper.getObservedVotes(ctx, assets.TerraDenom)
	require.Equal(t, 3, len(observedVotes))

	// Should now be two ObservedVotes
	keeper.deleteObservedVote(ctx, votes[0])
	observedVotes = keeper.getObservedVotes(ctx, assets.TerraDenom)
	require.Equal(t, 2, len(observedVotes))

	// Should now be zero ObservedVotes
	deleter := func(observedVote PriceVote) (stop bool) {
		keeper.deleteObservedVote(ctx, observedVote)
		return false
	}
	keeper.iterateObservedVotes(ctx, assets.TerraDenom, deleter)
	observedVotes = keeper.getObservedVotes(ctx, assets.TerraDenom)
	require.Equal(t, 0, len(observedVotes))
}
