package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParams(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, params.ValidateBasic())

	params = DefaultParams()
	params.TaxPolicy.RateMax = sdk.ZeroDec()
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.TaxPolicy.RateMin = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.TaxPolicy.Cap = sdk.Coin{Denom: "foo", Amount: sdk.NewInt(-1)}
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.TaxPolicy.ChangeRateMax = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.RewardPolicy.RateMax = sdk.ZeroDec()
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.RewardPolicy.ChangeRateMax = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.SeigniorageBurdenTarget = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.MiningIncrement = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.WindowShort = -1
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.WindowLong = -1
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.WindowProbation = -1
	require.Error(t, params.ValidateBasic())

	params = DefaultParams()
	params.RewardPolicy.RateMin = sdk.NewDec(-1)
	require.Error(t, params.ValidateBasic())

	require.NotNil(t, params.ParamSetPairs())
	require.NotNil(t, params.String())
}
