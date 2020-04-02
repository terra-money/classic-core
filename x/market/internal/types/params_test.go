package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_Validate(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, params.Validate())

	params = DefaultParams()
	params.BasePool = sdk.NewDec(-1)
	require.Error(t, params.Validate())

	params = DefaultParams()
	params.MinStabilitySpread = sdk.NewDec(-1)
	require.Error(t, params.Validate())
}
