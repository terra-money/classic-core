package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	err := p1.ValidateBasic()
	require.NoError(t, err)

	// invalid baes pool
	p1.BasePool = sdk.NewDec(-1)
	err = p1.ValidateBasic()
	require.Error(t, err)

	// invalid pool recovery period
	p2 := DefaultParams()
	p2.PoolRecoveryPeriod = 0
	err = p2.ValidateBasic()
	require.Error(t, err)

	// invalid min spread
	p3 := DefaultParams()
	p3.MinSpread = sdk.NewDecWithPrec(-1, 2)
	err = p3.ValidateBasic()
	require.Error(t, err)

	// invalid tobin tax
	p4 := DefaultParams()
	p4.TobinTax = sdk.NewDec(-1)
	err = p4.ValidateBasic()
	require.Error(t, err)

	// invalid illiquid tobin tax list
	p5 := DefaultParams()
	p5.IlliquidTobinTaxList = TobinTaxList{TobinTax{Denom: "foo", TaxRate: sdk.NewDec(-1)}}
	err = p5.ValidateBasic()
	require.Error(t, err)

	p6 := DefaultParams()
	require.NotNil(t, p6.ParamSetPairs())
	require.NotNil(t, p6.String())
}
