package test

import (
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFeeRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	taxAmount := sdk.NewInt(1000)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.KRWDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.GBPDenom, sdk.NewDec(100))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.CNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
		sdk.NewCoin(assets.SDRDenom, taxAmount),
		sdk.NewCoin(assets.KRWDenom, taxAmount),
		sdk.NewCoin(assets.GBPDenom, taxAmount),
		sdk.NewCoin(assets.CNYDenom, taxAmount),
	})

	// Get taxes
	taxProceedsInSDR := treasury.TaxRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.NewDec(1111), taxProceedsInSDR)
}

func TestSeigniorageRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	sAmt := sdk.NewInt(1000)
	lnasdrRate := sdk.NewDec(10)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, lnasdrRate)

	// Add seigniorage
	input.mintKeeper.AddSeigniorage(input.ctx, sAmt)

	// Get seigniorage rewards
	seigniorageProceeds := treasury.SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	miningRewardWeight := input.treasuryKeeper.GetRewardWeight(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, lnasdrRate.MulInt(sAmt).Mul(miningRewardWeight), seigniorageProceeds)
}

func TestMiningRewardsForEpoch(t *testing.T) {
	input := createTestInput(t)

	amt := sdk.NewInt(1000)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.KRWDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.GBPDenom, sdk.NewDec(100))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.CNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
		sdk.NewCoin(assets.SDRDenom, amt),
		sdk.NewCoin(assets.KRWDenom, amt),
		sdk.NewCoin(assets.GBPDenom, amt),
		sdk.NewCoin(assets.CNYDenom, amt),
	})

	// Add seigniorage
	input.mintKeeper.AddSeigniorage(input.ctx, amt)

	tProceeds := treasury.TaxRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	sProceeds := treasury.SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	mProceeds := treasury.MiningRewardForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))

	require.Equal(t, tProceeds.Add(sProceeds), mProceeds)
}

func TestSMR(t *testing.T) {
	input := createTestInput(t)
	amt := sdk.NewInt(1000)
	lnasdrRate := sdk.NewDec(10)

	// Set random prices
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, lnasdrRate)

	// Record tax proceeds
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{
		sdk.NewCoin(assets.SDRDenom, amt),
	})

	// Add seigniorage
	input.mintKeeper.AddSeigniorage(input.ctx, amt)

	tProceeds := treasury.TaxRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	sProceeds := treasury.SeigniorageRewardsForEpoch(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))

	actualSMR := treasury.SMR(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx))
	expectedSMR := sProceeds.Quo(tProceeds.Add(sProceeds))

	require.Equal(t, expectedSMR, actualSMR)
}

func TestUnitIndicator(t *testing.T) {
	input := createTestInput(t)

	lunaIssuance := sdk.NewInt(10000)
	err := input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.LunaDenom, lunaIssuance))
	require.Nil(t, err)

	// Just get an indicator to multiply the unit value by the expected rval.
	// the unit indicator function obviously should return the expected rval.
	actual := treasury.UnitLunaIndicator(input.ctx, input.treasuryKeeper, util.GetEpoch(input.ctx),
		func(_ sdk.Context, _ treasury.Keeper, _ sdk.Int) sdk.Dec {
			return sdk.NewDecFromInt(lunaIssuance.MulRaw(20))
		})

	require.Equal(t, sdk.NewDec(20), actual)
}

func linearFn(_ sdk.Context, _ treasury.Keeper, epoch sdk.Int) sdk.Dec {
	return sdk.NewDecFromInt(epoch)
}

func TestRollingAverageIndicator(t *testing.T) {
	input := createTestInput(t)

	// Case 1: at epoch 0 and averaging over 0 epochs
	rval := treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and averaging over negative epochs
	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.OneInt().Neg(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and averaging over 3, 4, 5 epochs; all should have the same rval
	input.ctx = input.ctx.WithBlockHeight(util.GetBlocksPerEpoch() * 3)
	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(4), linearFn)
	rval2 := treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(5), linearFn)
	rval3 := treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(6), linearFn)
	require.Equal(t, sdk.NewDecWithPrec(15, 1), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and averaging over 0 epochs
	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.ZeroInt(), linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5: at epoch 3 and averaging over 1 epoch
	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.OneInt(), linearFn)
	require.Equal(t, sdk.NewDec(3), rval)

	// Case 6: at epoch 500 and averaging over 300 epochs
	input.ctx = input.ctx.WithBlockHeight(util.GetBlocksPerEpoch() * 500)
	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), linearFn)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1), rval)

	// Test all of our reporting functions
	lunaAmt := int64(10000)
	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewInt64Coin(assets.LunaDenom, lunaAmt))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.OneDec())

	for i := int64(201); i <= 500; i++ {
		input.ctx = input.ctx.WithBlockHeight(util.GetBlocksPerEpoch() * i)
		input.treasuryKeeper.RecordTaxProceeds(input.ctx, sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(i))})
		input.mintKeeper.AddSeigniorage(input.ctx, sdk.NewInt(i))
		input.treasuryKeeper.SetRewardWeight(input.ctx, sdk.OneDec())
	}

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.TaxRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1), rval)

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.SeigniorageRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1), rval)

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.MiningRewardForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 1), rval)

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.TRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 5), rval)

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.SRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 5), rval)

	rval = treasury.RollingAverageIndicator(input.ctx, input.treasuryKeeper, sdk.NewInt(300), treasury.MRL)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 5), rval)
}
