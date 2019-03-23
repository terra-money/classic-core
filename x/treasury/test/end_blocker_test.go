package test

import (
	"terra/types/assets"
	"terra/x/treasury"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestEndBlockerTiming(t *testing.T) {
	input := createTestInput(t)

	// Housekeeping
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.SDRDenom, sdk.NewDec(1))

	// Test 1.
	resTags := treasury.EndBlocker(input.ctx, input.treasuryKeeper)
	pairs := resTags.ToKVPairs()
	require.Equal(t, len(pairs), 0)
}

func TestEndBlockerUpdateTaxPolicy(t *testing.T) {
}

func TestEndBlockerUpdateRewardPolicy(t *testing.T) {
}

func TestEndBlockerSettleClaims(t *testing.T) {
}
