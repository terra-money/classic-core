package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"testing"
)

func TestKeeper_ConvertRenewalFeeToTime(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, sdk.OneDec())
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, sdk.OneDec())
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroUSDDenom, sdk.OneDec().MulInt64(2))
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroMNTDenom, sdk.OneDec().QuoInt64(2))

	renewalInterval := input.NameserviceKeeper.RenewalInterval(input.Ctx)
	renewalFees := input.NameserviceKeeper.RenewalFees(input.Ctx)
	renewalFee := renewalFees.RenewalFeeForLength(3)

	// with default denom
	fee := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, renewalFee.Amount))
	extendedTime, err := input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.NoError(t, err)
	require.Equal(t, renewalInterval, extendedTime)

	// with same price denom
	fee = sdk.NewCoins(sdk.NewCoin(core.MicroKRWDenom, renewalFee.Amount))
	extendedTime, err = input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.NoError(t, err)
	require.Equal(t, renewalInterval, extendedTime)

	// with half price denom
	fee = sdk.NewCoins(sdk.NewCoin(core.MicroUSDDenom, renewalFee.Amount))
	extendedTime, err = input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.NoError(t, err)
	require.Equal(t, renewalInterval/2, extendedTime)

	// with double price denom
	fee = sdk.NewCoins(sdk.NewCoin(core.MicroMNTDenom, renewalFee.Amount))
	extendedTime, err = input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.NoError(t, err)
	require.Equal(t, renewalInterval*2, extendedTime)

	// with mixed prices
	fee = sdk.NewCoins(sdk.NewCoin(core.MicroMNTDenom, renewalFee.Amount), sdk.NewCoin(core.MicroUSDDenom, renewalFee.Amount))
	extendedTime, err = input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.NoError(t, err)
	require.Equal(t, renewalInterval*2+renewalInterval/2, extendedTime)

	// not registered prices
	fee = sdk.NewCoins(sdk.NewCoin("foo", renewalFee.Amount))
	_, err = input.NameserviceKeeper.ConvertRenewalFeeToTime(input.Ctx, fee, 3)
	require.Error(t, err)
}
