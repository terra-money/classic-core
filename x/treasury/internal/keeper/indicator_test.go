package keeper

import (
	"fmt"
	"testing"

	core "github.com/terra-project/core/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestFeeRewardsForEpoch(t *testing.T) {
	input := CreateTestInput(t)

	taxAmount := sdk.NewInt(1000).MulRaw(core.MicroUnit)

	// Set random prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.NewDec(1))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, sdk.NewDec(10))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroGBPDenom, sdk.NewDec(100))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroCNYDenom, sdk.NewDec(1000))

	// Record tax proceeds
	input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, sdk.Coins{
		sdk.NewCoin(core.MicroSDRDenom, taxAmount),
		sdk.NewCoin(core.MicroKRWDenom, taxAmount),
		sdk.NewCoin(core.MicroGBPDenom, taxAmount),
		sdk.NewCoin(core.MicroCNYDenom, taxAmount),
	})

	// Get taxes
	taxProceedsInSDR := TaxRewardsForEpoch(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx))
	require.Equal(t, sdk.NewDec(1111).MulInt64(core.MicroUnit), taxProceedsInSDR)
}

func TestSeigniorageRewardsForEpoch(t *testing.T) {
	input := CreateTestInput(t)

	sAmt := sdk.NewInt(1000)
	lnasdrRate := sdk.NewDec(10)

	// Add seigniorage
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sAmt)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordHistoricalIssuance(input.Ctx)

	// Set random prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lnasdrRate)
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch)

	// Add seigniorage
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	// Get seigniorage rewards
	seigniorageProceeds := SeigniorageRewardsForEpoch(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx))
	miningRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx, core.GetEpoch(input.Ctx))
	require.Equal(t, lnasdrRate.MulInt(sAmt).Mul(miningRewardWeight), seigniorageProceeds)
}

func TestMiningRewardsForEpoch(t *testing.T) {
	input := CreateTestInput(t)

	amt := sdk.NewInt(1000).MulRaw(core.MicroUnit)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, amt)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordHistoricalIssuance(input.Ctx)

	// Set random prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, sdk.NewDec(1))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.NewDec(10))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroGBPDenom, sdk.NewDec(100))
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroCNYDenom, sdk.NewDec(1000))

	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch)

	// Record tax proceeds
	input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, sdk.Coins{
		sdk.NewCoin(core.MicroSDRDenom, amt),
		sdk.NewCoin(core.MicroKRWDenom, amt),
		sdk.NewCoin(core.MicroGBPDenom, amt),
		sdk.NewCoin(core.MicroCNYDenom, amt),
	})

	// Add seigniorage
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	tProceeds := TaxRewardsForEpoch(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx))
	sProceeds := SeigniorageRewardsForEpoch(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx))
	mProceeds := MiningRewardForEpoch(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx))

	require.Equal(t, tProceeds.Add(sProceeds), mProceeds)
}

func TestUnitIndicator(t *testing.T) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(100)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	res := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	fmt.Println(res)
	require.True(t, res.IsOK())
	res = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t, res.IsOK())
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	lunaTotalBondedAmount := input.StakingKeeper.TotalBondedTokens(input.Ctx)

	// Just get an indicator to multiply the unit value by the expected rval.
	// the unit indicator function obviously should return the expected rval.
	actual := UnitLunaIndicator(input.Ctx, input.TreasuryKeeper, core.GetEpoch(input.Ctx),
		func(_ sdk.Context, _ Keeper, _ int64) sdk.Dec {
			return sdk.NewDecFromInt(lunaTotalBondedAmount.MulRaw(20))
		})

	require.Equal(t, sdk.NewDec(20), actual)
}

func linearFn(_ sdk.Context, _ Keeper, epoch int64) sdk.Dec {
	return sdk.NewDec(epoch)
}

func TestSumIndicator(t *testing.T) {
	input := CreateTestInput(t)

	// Case 1: at epoch 0 and summing over 0 epochs
	rval := SumIndicator(input.Ctx, input.TreasuryKeeper, 0, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and summing over negative epochs
	rval = SumIndicator(input.Ctx, input.TreasuryKeeper, -1, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and summing over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 3)
	rval = SumIndicator(input.Ctx, input.TreasuryKeeper, 4, linearFn)
	rval2 := SumIndicator(input.Ctx, input.TreasuryKeeper, 5, linearFn)
	rval3 := SumIndicator(input.Ctx, input.TreasuryKeeper, 6, linearFn)
	require.Equal(t, sdk.NewDec(6), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and summing over 0 epochs
	rval = SumIndicator(input.Ctx, input.TreasuryKeeper, 0, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5. Sum up to 10
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 10)
	rval = SumIndicator(input.Ctx, input.TreasuryKeeper, 10, linearFn)
	require.Equal(t, sdk.NewDec(55), rval)
}

func TestRollingAverageIndicator(t *testing.T) {
	input := CreateTestInput(t)
	sh := staking.NewHandler(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(1)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	res := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.True(t, res.IsOK())
	res = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t, res.IsOK())
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	// Case 1: at epoch 0 and averaging over 0 epochs
	rval := RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 0, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and averaging over negative epochs
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, -1, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and averaging over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 3)
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 4, linearFn)
	rval2 := RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 5, linearFn)
	rval3 := RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 6, linearFn)
	require.Equal(t, sdk.NewDecWithPrec(15, 1), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and averaging over 0 epochs
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 0, linearFn)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5: at epoch 3 and averaging over 1 epoch
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 1, linearFn)
	require.Equal(t, sdk.NewDec(3), rval)

	// Case 6: at epoch 500 and averaging over 300 epochs
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 500)
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, linearFn)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1), rval)

	// Test all of our reporting functions
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	// set initial supply
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 200)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100000000*core.MicroUnit))))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordHistoricalIssuance(input.Ctx)

	for i := int64(201); i <= 500; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * i)
		input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, sdk.Coins{sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(i).MulRaw(core.MicroUnit))})
		input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.OneDec())

		supply = supply.SetTotal(supply.GetTotal().Sub(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(i).MulRaw(core.MicroUnit)))))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)
		input.TreasuryKeeper.RecordHistoricalIssuance(input.Ctx)
	}

	totalBondedTokens := sdk.NewDecFromInt(input.StakingKeeper.TotalBondedTokens(input.Ctx))
	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, TaxRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(core.MicroUnit), rval)

	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, SeigniorageRewardsForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(core.MicroUnit), rval)

	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, MiningRewardForEpoch)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 1).MulInt64(core.MicroUnit), rval)

	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, TRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(core.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.Mul(sdk.NewDec(1000000)).TruncateInt())

	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, SRL)
	require.Equal(t, sdk.NewDecWithPrec(3505, 1).MulInt64(core.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.MulTruncate(sdk.NewDec(1000000)).TruncateInt())

	rval = RollingAverageIndicator(input.Ctx, input.TreasuryKeeper, 300, MRL)
	require.Equal(t, sdk.NewDecWithPrec(3505*2, 1).MulInt64(core.MicroUnit).Quo(totalBondedTokens).Mul(sdk.NewDec(1000000)).TruncateInt(), rval.MulTruncate(sdk.NewDec(1000000)).TruncateInt())
}
