package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Test_DenomList(t *testing.T) {
	denoms := types.DenomList{
		types.Denom{
			Name:     "denom1",
			TobinTax: sdk.NewDec(100),
		},
		types.Denom{
			Name:     "denom2",
			TobinTax: sdk.NewDec(200),
		},
		types.Denom{
			Name:     "denom3",
			TobinTax: sdk.NewDec(300),
		},
	}

	require.False(t, denoms[0].Equal(&denoms[1]))
	require.True(t, denoms[0].Equal(&denoms[0]))
	require.Equal(t, "name: denom1\ntobin_tax: \"100.000000000000000000\"\n", denoms[0].String())
	require.Equal(t, "name: denom2\ntobin_tax: \"200.000000000000000000\"\n", denoms[1].String())
	require.Equal(t, "name: denom3\ntobin_tax: \"300.000000000000000000\"\n", denoms[2].String())
	require.Equal(t, "name: denom1\ntobin_tax: \"100.000000000000000000\"\n\nname: denom2\ntobin_tax: \"200.000000000000000000\"\n\nname: denom3\ntobin_tax: \"300.000000000000000000\"", denoms.String())
}
