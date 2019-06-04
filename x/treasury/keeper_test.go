package treasury

import (
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRewardWeight(t *testing.T) {
	input := createTestInput(t)

	// See that we can get and set reward weights
	blocksPerEpoch := util.BlocksPerEpoch
	for i := int64(0); i < 10; i++ {
		input.ctx = input.ctx.WithBlockHeight(i * blocksPerEpoch)

		input.treasuryKeeper.SetRewardWeight(input.ctx, sdk.NewDecWithPrec(i, 2))
	}

	for i := int64(0); i < 10; i++ {
		input.ctx = input.ctx.WithBlockHeight(i * blocksPerEpoch)

		require.Equal(t, sdk.NewDecWithPrec(i, 2), input.treasuryKeeper.GetRewardWeight(input.ctx, sdk.NewInt(i)))
	}
}

func TestTax(t *testing.T) {
	input := createTestInput(t)

	// Set & get tax rate
	testRate := sdk.NewDecWithPrec(2, 3)
	input.treasuryKeeper.SetTaxRate(input.ctx, testRate)
	curRate := input.treasuryKeeper.GetTaxRate(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, curRate, testRate)

	// Vicariously set tax caps & test
	params := DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, params)
	sdrCap := params.TaxPolicy.Cap

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, sdk.NewDec(100))

	readSdrCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.MicroSDRDenom)
	cnyCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.MicroCNYDenom)
	krwCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.MicroKRWDenom)

	require.Equal(t, sdrCap.Amount, readSdrCap)
	require.Equal(t, sdrCap.Amount.MulRaw(10), cnyCap)
	require.Equal(t, sdrCap.Amount.MulRaw(100), krwCap)
}

func TestParams(t *testing.T) {
	input := createTestInput(t)

	defaultParams := DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, defaultParams)

	retrievedParams := input.treasuryKeeper.GetParams(input.ctx)
	require.Equal(t, defaultParams, retrievedParams)
}
