package treasury

import (
	"terra/types/assets"
	"terra/types/util"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestConstraintClamp(t *testing.T) {
	input := createTestInput(t)

	taxPolicy := input.treasuryKeeper.GetParams(input.ctx).TaxPolicy
	prevRate := input.treasuryKeeper.GetTaxRate(input.ctx, util.GetEpoch(input.ctx))

	// Case 1: try to update delta > maxUpdateRate
	newRate := prevRate.Add(taxPolicy.ChangeRateMax.MulInt64(2))
	clampedRate := taxPolicy.Clamp(prevRate, newRate)
	require.Equal(t, prevRate.Add(taxPolicy.ChangeRateMax), clampedRate)

	// Case 2: try to update delta > maxUpdateRate in other direction
	newRate = prevRate.Sub(taxPolicy.ChangeRateMax.MulInt64(2))
	clampedRate = taxPolicy.Clamp(prevRate, newRate)
	require.Equal(t, prevRate.Sub(taxPolicy.ChangeRateMax), clampedRate)

	// Case 3: try to update the new rate > maxRate
	prevRate = taxPolicy.RateMax
	newRate = taxPolicy.RateMax.Add(sdk.NewDecWithPrec(1, 3))
	clampedRate = taxPolicy.Clamp(prevRate, newRate)
	require.Equal(t, taxPolicy.RateMax, clampedRate)

	// Case 4: try to update the new rate < minRate
	prevRate = taxPolicy.RateMin
	newRate = taxPolicy.RateMin.Sub(sdk.NewDecWithPrec(1, 3))
	clampedRate = taxPolicy.Clamp(prevRate, newRate)
	require.Equal(t, taxPolicy.RateMin, clampedRate)
}

func TestConstraintCap(t *testing.T) {
	input := createTestInput(t)
	taxPolicy := input.treasuryKeeper.GetParams(input.ctx).TaxPolicy

	// Set prices for test assets first
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.KRWDenom, sdk.NewDec(1000))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.GBPDenom, sdk.NewDec(1))

	// Check that SDR tax cap has been set
	require.Equal(t, taxPolicy.Cap.Amount, input.treasuryKeeper.GetTaxCap(input.ctx, assets.SDRDenom))
	require.Equal(t, sdk.NewInt(100), input.treasuryKeeper.GetTaxCap(input.ctx, assets.KRWDenom))
	require.Equal(t, sdk.NewInt(1), input.treasuryKeeper.GetTaxCap(input.ctx, assets.GBPDenom))
}
