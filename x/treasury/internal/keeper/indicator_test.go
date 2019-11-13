package keeper

import (
	"testing"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFeeRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

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

	// Update Indicators
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Get Tax Rewards per Luna (TRL) & Tax Rawards (TR)
	TRL := input.TreasuryKeeper.GetTRL(input.Ctx, core.GetEpoch(input.Ctx))
	TR := TRL.MulInt(input.StakingKeeper.TotalBondedTokens(input.Ctx))
	require.Equal(t, sdk.NewDec(1111).MulInt64(core.MicroUnit), TR)
}

func TestSeigniorageRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

	sAmt := sdk.NewInt(1000)
	lnasdrRate := sdk.NewDec(10)

	// Add seigniorage
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sAmt)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	// Set random prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lnasdrRate)
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch)

	// Add seigniorage
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	// Update Indicators
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Get seigniorage rewards (SR)
	SR := input.TreasuryKeeper.GetSR(input.Ctx, core.GetEpoch(input.Ctx))
	miningRewardWeight := input.TreasuryKeeper.GetRewardWeight(input.Ctx)
	require.Equal(t, lnasdrRate.MulInt(sAmt).Mul(miningRewardWeight), SR)
}

func TestMiningRewardsForEpoch(t *testing.T) {
	input, _ := setupValidators(t)

	amt := sdk.NewInt(1000).MulRaw(core.MicroUnit)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, amt)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

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

	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	epoch := core.GetEpoch(input.Ctx)

	tProceeds := input.TreasuryKeeper.GetTRL(input.Ctx, epoch).MulInt(input.StakingKeeper.TotalBondedTokens(input.Ctx))
	sProceeds := input.TreasuryKeeper.GetSR(input.Ctx, epoch)
	mProceeds := input.TreasuryKeeper.GetMR(input.Ctx, epoch)

	require.Equal(t, tProceeds.Add(sProceeds), mProceeds)
}

func TestLoadIndicatorByEpoch(t *testing.T) {
	input := CreateTestInput(t)

	MRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
	}

	for epoch, MR := range MRArr {
		input.TreasuryKeeper.SetMR(input.Ctx, int64(epoch), MR)
	}

	SRArr := []sdk.Dec{
		sdk.NewDec(10),
		sdk.NewDec(20),
		sdk.NewDec(30),
		sdk.NewDec(40),
	}

	for epoch, SR := range SRArr {
		input.TreasuryKeeper.SetSR(input.Ctx, int64(epoch), SR)
	}

	TRLArr := []sdk.Dec{
		sdk.NewDecWithPrec(1, 5),
		sdk.NewDecWithPrec(2, 4),
		sdk.NewDecWithPrec(3, 3),
		sdk.NewDecWithPrec(4, 2),
	}

	for epoch, TRL := range TRLArr {
		input.TreasuryKeeper.SetTRL(input.Ctx, int64(epoch), TRL)
	}

	for epoch := int64(0); epoch < 4; epoch++ {
		require.Equal(t, MRArr[epoch], input.TreasuryKeeper.loadIndicatorByEpoch(input.Ctx, types.MRKey, epoch))
		require.Equal(t, SRArr[epoch], input.TreasuryKeeper.loadIndicatorByEpoch(input.Ctx, types.SRKey, epoch))
		require.Equal(t, TRLArr[epoch], input.TreasuryKeeper.loadIndicatorByEpoch(input.Ctx, types.TRLKey, epoch))
	}
}

func linearFn(_ sdk.Context, _ Keeper, epoch int64) sdk.Dec {
	return sdk.NewDec(epoch)
}

func TestSumIndicator(t *testing.T) {
	input := CreateTestInput(t)

	MRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
		sdk.NewDec(500),
		sdk.NewDec(600),
	}

	for epoch, MR := range MRArr {
		input.TreasuryKeeper.SetMR(input.Ctx, int64(epoch), MR)
	}

	// Case 1: at epoch 0 and summing over 0 epochs
	rval := input.TreasuryKeeper.sumIndicator(input.Ctx, 0, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and summing over negative epochs
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, -1, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and summing over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 3)
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 4, types.MRKey)
	rval2 := input.TreasuryKeeper.sumIndicator(input.Ctx, 5, types.MRKey)
	rval3 := input.TreasuryKeeper.sumIndicator(input.Ctx, 6, types.MRKey)
	require.Equal(t, sdk.NewDec(1000), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and summing over 0 epochs
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 0, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5. Sum up to 6
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 5)
	rval = input.TreasuryKeeper.sumIndicator(input.Ctx, 6, types.MRKey)
	require.Equal(t, sdk.NewDec(2100), rval)
}

func TestRollingAverageIndicator(t *testing.T) {
	input := CreateTestInput(t)
	MRArr := []sdk.Dec{
		sdk.NewDec(100),
		sdk.NewDec(200),
		sdk.NewDec(300),
		sdk.NewDec(400),
	}

	for epoch, MR := range MRArr {
		input.TreasuryKeeper.SetMR(input.Ctx, int64(epoch), MR)
	}

	// Case 1: at epoch 0 and averaging over 0 epochs
	rval := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 0, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 2: at epoch 0 and averaging over negative epochs
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, -1, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 3: at epoch 3 and averaging over 3, 4, 5 epochs; all should have the same rval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch * 3)
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 4, types.MRKey)
	rval2 := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 5, types.MRKey)
	rval3 := input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 6, types.MRKey)
	require.Equal(t, sdk.NewDec(250), rval)
	require.Equal(t, rval, rval2)
	require.Equal(t, rval2, rval3)

	// Case 4: at epoch 3 and averaging over 0 epochs
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 0, types.MRKey)
	require.Equal(t, sdk.ZeroDec(), rval)

	// Case 5: at epoch 3 and averaging over 1 epoch
	rval = input.TreasuryKeeper.rollingAverageIndicator(input.Ctx, 1, types.MRKey)
	require.Equal(t, sdk.NewDec(400), rval)
}
