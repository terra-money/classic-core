package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"
)

func TestRewardWeight(t *testing.T) {
	input := CreateTestInput(t)

	// See that we can get and set reward weights
	for i := int64(0); i < 10; i++ {
		input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.NewDecWithPrec(i, 2))
		require.Equal(t, sdk.NewDecWithPrec(i, 2), input.TreasuryKeeper.GetRewardWeight(input.Ctx))
	}
}

func TestTaxRate(t *testing.T) {
	input := CreateTestInput(t)

	// See that we can get and set tax rate
	for i := int64(0); i < 10; i++ {
		input.TreasuryKeeper.SetTaxRate(input.Ctx, sdk.NewDecWithPrec(i, 2))
		require.Equal(t, sdk.NewDecWithPrec(i, 2), input.TreasuryKeeper.GetTaxRate(input.Ctx))
	}
}

func TestTaxCap(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroCNYDenom, sdk.NewInt(i))
		require.Equal(t, sdk.NewInt(i), input.TreasuryKeeper.GetTaxCap(input.Ctx, core.MicroCNYDenom))
	}
}

func TestIterateTaxCap(t *testing.T) {
	input := CreateTestInput(t)

	cnyCap := sdk.NewInt(123)
	usdCap := sdk.NewInt(13)
	krwCap := sdk.NewInt(1300)
	input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroCNYDenom, cnyCap)
	input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroUSDDenom, usdCap)
	input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroKRWDenom, krwCap)

	input.TreasuryKeeper.IterateTaxCap(input.Ctx, func(denom string, taxCap sdk.Int) bool {
		switch denom {
		case core.MicroCNYDenom:
			require.Equal(t, cnyCap, taxCap)
		case core.MicroUSDDenom:
			require.Equal(t, usdCap, taxCap)
		case core.MicroKRWDenom:
			require.Equal(t, krwCap, taxCap)
		}

		return false
	})

}

func TestTaxProceeds(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		proceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(100+i)))
		input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, proceeds)
		input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, proceeds)
		input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, proceeds)

		require.Equal(t, proceeds.Add(proceeds).Add(proceeds), input.TreasuryKeeper.PeekEpochTaxProceeds(input.Ctx))

		input.TreasuryKeeper.SetEpochTaxProceeds(input.Ctx, sdk.Coins{})
		require.Equal(t, sdk.Coins{}, input.TreasuryKeeper.PeekEpochTaxProceeds(input.Ctx))
	}
}

func TestMicroLunaIssuance(t *testing.T) {
	input := CreateTestInput(t)

	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	// See that we can get and set luna issuance
	blocksPerEpoch := core.BlocksPerWeek
	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(i))))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)
		input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

		require.Equal(t, sdk.NewInt(i), input.TreasuryKeeper.GetEpochInitialIssuance(input.Ctx).AmountOf(core.MicroLunaDenom))
	}
}

func TestPeekEpochSeigniorage(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * core.BlocksPerWeek)
		supply := input.SupplyKeeper.GetSupply(input.Ctx)

		preIssuance := sdk.NewInt(rand.Int63() + 1)
		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, preIssuance)))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)
		input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

		nowIssuance := sdk.NewInt(rand.Int63() + 1)
		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, nowIssuance)))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)

		targetSeigniorage := preIssuance.Sub(nowIssuance)
		if targetSeigniorage.IsNegative() {
			targetSeigniorage = sdk.ZeroInt()
		}

		require.Equal(t, targetSeigniorage, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx))
	}
}

func TestCumulatedHeight(t *testing.T) {
	input := CreateTestInput(t)

	// See that we can get and set reward weights
	for i := int64(0); i < 10; i++ {
		input.TreasuryKeeper.SetCumulatedHeight(input.Ctx, i*100)
		require.Equal(t, i*100, input.TreasuryKeeper.GetCumulatedHeight(input.Ctx))
	}
}

func TestIndicatorGetterSetter(t *testing.T) {
	input := CreateTestInput(t)

	for e := int64(0); e < 10; e++ {
		randomVal := sdk.NewDec(rand.Int63() + 1)
		input.TreasuryKeeper.SetTR(input.Ctx, e, randomVal)
		require.Equal(t, randomVal, input.TreasuryKeeper.GetTR(input.Ctx, e))
		input.TreasuryKeeper.SetSR(input.Ctx, e, randomVal)
		require.Equal(t, randomVal, input.TreasuryKeeper.GetSR(input.Ctx, e))
		input.TreasuryKeeper.SetTSL(input.Ctx, e, randomVal.TruncateInt())
		require.Equal(t, randomVal.TruncateInt(), input.TreasuryKeeper.GetTSL(input.Ctx, e))
	}

	input.TreasuryKeeper.ClearTRs(input.Ctx)
	input.TreasuryKeeper.ClearSRs(input.Ctx)
	input.TreasuryKeeper.ClearTSLs(input.Ctx)

	for e := int64(0); e < 10; e++ {
		require.Equal(t, sdk.ZeroDec(), input.TreasuryKeeper.GetTR(input.Ctx, e))
		require.Equal(t, sdk.ZeroDec(), input.TreasuryKeeper.GetSR(input.Ctx, e))
		require.Equal(t, sdk.ZeroInt(), input.TreasuryKeeper.GetTSL(input.Ctx, e))
	}
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	defaultParams := types.DefaultParams()
	input.TreasuryKeeper.SetParams(input.Ctx, defaultParams)

	retrievedParams := input.TreasuryKeeper.GetParams(input.Ctx)
	require.Equal(t, defaultParams, retrievedParams)
}
