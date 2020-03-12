package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"testing"
)

func TestRenewalFees_RenewalFeeForLength(t *testing.T) {
	feeFor3 := sdk.NewInt64Coin(core.MicroSDRDenom, 100)
	feeFor4 := sdk.NewInt64Coin(core.MicroSDRDenom, 5)
	feeFor5 := sdk.NewInt64Coin(core.MicroSDRDenom, 2)
	feeFor10 := sdk.NewInt64Coin(core.MicroSDRDenom, 1)

	renewalFees := RenewalFees{
		{4, feeFor4},
		{5, feeFor5},
		{3, feeFor3},
		{10, feeFor10},
	}

	require.Equal(t, feeFor3, renewalFees.RenewalFeeForLength(3))
	require.Equal(t, feeFor4, renewalFees.RenewalFeeForLength(4))
	require.Equal(t, feeFor5, renewalFees.RenewalFeeForLength(5))
	require.Equal(t, feeFor10, renewalFees.RenewalFeeForLength(6))
	require.Equal(t, feeFor10, renewalFees.RenewalFeeForLength(10))
	require.Equal(t, feeFor10, renewalFees.RenewalFeeForLength(11))
}
