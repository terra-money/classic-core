package treasury

import (
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFeeRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	taxAmount := sdk.NewInt(1000).MulRaw(assets.MicroUnit)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroGBPDenom, sdk.NewDec(100))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
		sdk.NewCoin(assets.MicroSDRDenom, taxAmount),
		sdk.NewCoin(assets.MicroKRWDenom, taxAmount),
		sdk.NewCoin(assets.MicroGBPDenom, taxAmount),
		sdk.NewCoin(assets.MicroCNYDenom, taxAmount),
	})

	// Get taxes
	taxProceedsInSDR := TaxRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.NewDec(1111).MulInt64(assets.MicroUnit), taxProceedsInSDR)
}

func TestSeigniorageRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	sAmt := sdk.NewInt(1000)
	lnasdrRate := sdk.NewDec(10)

	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, sAmt))

	SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, lnasdrRate)

	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch)

	// Add seigniorage
	input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, sAmt))

	// Get seigniorage rewards
	seigniorageProceeds := SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	miningRewardWeight := input.treasuryKeeper.GetRewardWeight(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, lnasdrRate.MulInt(sAmt).Mul(miningRewardWeight), seigniorageProceeds)
}

func TestMiningRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	amt := sdk.NewInt(1000).MulRaw(assets.MicroUnit)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroGBPDenom, sdk.NewDec(100))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
		sdk.NewCoin(assets.MicroSDRDenom, amt),
		sdk.NewCoin(assets.MicroKRWDenom, amt),
		sdk.NewCoin(assets.MicroGBPDenom, amt),
		sdk.NewCoin(assets.MicroCNYDenom, amt),
	})

	// Add seigniorage
	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, amt))

	tProceeds := TaxRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	sProceeds := SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	mProceeds := MiningRewardForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))

	require.Equal(t, tProceeds.Add(sProceeds), mProceeds)
}

func TestUnitIndicator(t *testing.T) {
	input := createTestInput(t)

	lunaTotalBondedAmount := input.treasuryKeeper.valset.TotalBondedTokens(input.ctx)

	// Just get an indicator to multiply the unit value by the expected rval.
	// the unit indicator function obviously should return the expected rval.
	actual := UnitLunaIndicator(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx),
		func(_ sdk.Context, _ Keeper, _ sdk.Int) sdk.Dec {
			return sdk.NewDecFromInt(lunaTotalBondedAmount.MulRaw(20))
		})

	require.Equal(t, sdk.NewDec(20), actual)
}

func linearFn(_ sdk.Context, _ Keeper, epoch sdk.Int) sdk.Dec {
	return sdk.NewDecFromInt(epoch)
}

func TestSumIndicator(t *testing.T) {
	input := createTestInput(t)

	// Case 1: at epoch 0 and summing over 0 epochs
	rval := SumIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and summing over negative epochs
	rval = SumIndicator(input.ctx, input.treasuryKeeper, sdk.OneInt().Neg(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and summing over 3, 4, 5 epochs; all should have the same rval
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * 3)
	rval = SumIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(4), linearFn)
	rval2 := SumIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(5), linearFn)
	rval3 := SumIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(6), linearFn)
	require.Equal(t, sdk.NewDec(6), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and summing over 0 epochs
	rval = SumIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5. Sum up to 10
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * 10)
	rval = SumIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(10), linearFn)
	require.Equal(t, sdk.NewDec(55), rval)
}

func TestRollingAverageIndicator(t *testing.T) {
	input := createTestInput(t)

	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(100000000*assets.MicroUnit)))

	// Case 1: at epoch 0 and averaging over 0 epochs
	rval := RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and averaging over negative epochs
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.OneInt().Neg(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and averaging over 3, 4, 5 epochs; all should have the same rval
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * 3)
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(4), linearFn)
	rval2 := RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(5), linearFn)
	rval3 := RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(6), linearFn)
	require.Equal(t, sdk.NewDecWithPrec(15, 1), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and averaging over 0 epochs
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5: at epoch 3 and averaging over 1 epoch
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.OneInt(), linearFn)
	require.Equal(t, sdk.NewDec(3), rval)

	// Case 6: at epoch 500 and averaging over 300 epochs
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * 500)
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), linearFn)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1), rval)

	// Test all of our reporting functions
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.OneDec())
	input.mintKeeper.PeekEpochSeigniorage(input.ctx, sdk.ZeroInt())

	for i := int64(201); i <= 500; i++ {
		input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * i)
		input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(i).MulRaw(assets.MicroUnit))})
		input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(i).MulRaw(assets.MicroUnit)))

		input.treasuryKeeper.SetRewardWeight(input.ctx, sdk.OneDec())
	}

	totalBondedTokens := sdk.NewDecFromInt(input.treasuryKeeper.valset.TotalBondedTokens(input.ctx))
	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), TaxRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(assets.MicroUnit), rval)

	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), SeigniorageRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(assets.MicroUnit), rval)

	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), MiningRewardForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 1).MulInt64(assets.MicroUnit), rval)

	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), TRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(assets.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.Mul(sdk.NewDec(1000000)).TruncateInt())

	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), SRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(assets.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.MulTruncate(sdk.NewDec(1000000)).TruncateInt())

	rval = RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), MRL)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 1).MulInt64(assets.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.MulTruncate(sdk.NewDec(1000000)).TruncateInt())
}
