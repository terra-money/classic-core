package distribution

import (
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/mint"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
)

// set the proposer for determining distribution during endblock
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, dk keeper.Keeper, fck auth.FeeCollectionKeeper, mtk market.Keeper, mk mint.Keeper) {

	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		alignFee(ctx, fck, mtk, mk)
		previousProposer := dk.GetPreviousProposerConsAddr(ctx)
		dk.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	dk.SetPreviousProposerConsAddr(ctx, consAddr)

}

func alignFee(ctx sdk.Context, fck auth.FeeCollectionKeeper, mtk market.Keeper, mk mint.Keeper) {
	// Swap feepool to SDR
	alignedFeePool := sdk.NewCoins()
	feepool := fck.GetCollectedFees(ctx)
	for _, coin := range feepool {
		if coin.Denom == assets.MicroSDRDenom {
			alignedFeePool = alignedFeePool.Add(sdk.NewCoins(coin))
			continue
		}

		retCoin, _, err := mtk.GetSwapCoin(ctx, coin, assets.MicroSDRDenom, true)
		if err != nil {
			alignedFeePool = alignedFeePool.Add(sdk.NewCoins(coin))
			continue
		}

		err = mk.InternalBurn(ctx, coin)
		if err != nil {
			panic(err)
		}

		err = mk.InternalMint(ctx, retCoin)
		if err != nil {
			panic(err)
		}

		alignedFeePool = alignedFeePool.Add(sdk.NewCoins(retCoin))
	}

	fck.ClearCollectedFees(ctx)
	fck.AddCollectedFees(ctx, alignedFeePool)
}
