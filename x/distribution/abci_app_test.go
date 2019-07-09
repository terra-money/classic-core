package distribution

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/types/assets"
)

func TestAlignFee(t *testing.T) {
	input := createTestInput(t)

	// Case 1: all fees should be aligned to SDR
	initialFee := sdk.NewCoins(
		sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(assets.MicroKRWDenom, sdk.NewInt(1000000)),
	)

	input.feeCollectionKeeper.AddCollectedFees(input.ctx, initialFee)

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.NewDecWithPrec(5, 1)) // 0.5
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, sdk.NewDecWithPrec(1, 1)) // 0.1

	alignFee(input.ctx, input.feeCollectionKeeper, input.marketKeeper, input.mintKeeper)

	alignedFeePool := input.feeCollectionKeeper.GetCollectedFees(input.ctx)
	require.Equal(t, 1, alignedFeePool.Len())
	require.Equal(t, assets.MicroSDRDenom, alignedFeePool[0].Denom)

	// Case 2: denom with no effective price, will not be aligned to SDR
	initialFee = sdk.NewCoins(
		sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(assets.MicroUSDDenom, sdk.NewInt(1000000)),
	)

	input.feeCollectionKeeper.ClearCollectedFees(input.ctx)
	input.feeCollectionKeeper.AddCollectedFees(input.ctx, initialFee)

	alignFee(input.ctx, input.feeCollectionKeeper, input.marketKeeper, input.mintKeeper)
	alignedFeePool = input.feeCollectionKeeper.GetCollectedFees(input.ctx)
	require.Equal(t, 2, alignedFeePool.Len())
}
