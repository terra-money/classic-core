package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"testing"
)

func TestParams_Validate(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, params.Validate())

	// invalid bid period
	params = DefaultParams()
	params.BidPeriod = 1
	require.Error(t, params.Validate())

	// invalid reveal period
	params = DefaultParams()
	params.RevealPeriod = 1
	require.Error(t, params.Validate())

	// invalid renew interval period
	params = DefaultParams()
	params.RenewalInterval = 1
	require.Error(t, params.Validate())

	// invalid min deposit
	params = DefaultParams()
	params.MinDeposit = sdk.Coin{}
	require.Error(t, params.Validate())

	// invalid root name
	params = DefaultParams()
	params.RootName = ""
	require.Error(t, params.Validate())

	// invalid renewal fees
	// invalid coin, invalid length, duplicated renewal fees
	params = DefaultParams()
	params.RenewalFees = RenewalFees{{1, sdk.Coin{Denom: core.MicroSDRDenom, Amount: sdk.NewInt(-1)}}}
	require.Error(t, params.Validate())

	params.RenewalFees = RenewalFees{{-1, sdk.Coin{Denom: core.MicroSDRDenom, Amount: sdk.NewInt(1)}}}
	require.Error(t, params.Validate())

	params.RenewalFees = RenewalFees{
		{1, sdk.Coin{Denom: core.MicroSDRDenom, Amount: sdk.NewInt(1)}},
		{1, sdk.Coin{Denom: core.MicroSDRDenom, Amount: sdk.NewInt(1)}},
	}
	require.Error(t, params.Validate())

}
