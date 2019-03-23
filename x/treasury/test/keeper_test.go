package test

import (
	"terra/types"
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRewardWeight(t *testing.T) {
	input := createTestInput(t)

	// See that we can get and set reward weights
	blocksPerEpoch := util.GetBlocksPerEpoch()
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
	curRate := input.treasuryKeeper.GetTaxRate(input.ctx)
	require.Equal(t, curRate, testRate)

	// Vicariously set tax caps & test
	params := treasury.DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, params)
	sdrCap := params.TaxPolicy.Cap

	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(1))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.CNYDenom, sdk.NewDec(10))
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.KRWDenom, sdk.NewDec(100))

	readSdrCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.SDRDenom)
	cnyCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.CNYDenom)
	krwCap := input.treasuryKeeper.GetTaxCap(input.ctx, assets.KRWDenom)

	require.Equal(t, sdrCap.Amount, readSdrCap)
	require.Equal(t, sdrCap.Amount.MulRaw(10), cnyCap)
	require.Equal(t, sdrCap.Amount.MulRaw(100), krwCap)
}

func TestClaim(t *testing.T) {
	input := createTestInput(t)

	for i := 0; i < 99; i++ {
		oracleClaim := types.NewClaim(
			types.OracleClaimClass, sdk.OneInt(), addrs[i%3],
		)
		input.treasuryKeeper.AddClaim(input.ctx, oracleClaim)

		budgetClaim := types.NewClaim(
			types.OracleClaimClass, sdk.OneInt(), addrs[i%3],
		)
		input.treasuryKeeper.AddClaim(input.ctx, budgetClaim)
	}

	// There should only be 3 unique claims, for each of the three addresses.
	// Each claim should have coalesced its weight to 33.
	counter := 0
	input.treasuryKeeper.IterateClaims(input.ctx, func(claim types.Claim) bool {
		counter++
		require.Equal(t, int64(66), claim.Weight.Int64())
		return false
	})

	require.Equal(t, 3, counter)
}

func TestParams(t *testing.T) {
	input := createTestInput(t)

	defaultParams := treasury.DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, defaultParams)

	retrievedParams := input.treasuryKeeper.GetParams(input.ctx)
	require.Equal(t, defaultParams, retrievedParams)
}
