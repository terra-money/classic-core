package pay

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIssuance(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	// New context. There should be no price.
	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

	terraPrice := sdk.NewDecWithPrec(166, 2)
	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(terraPrice) && err == nil)
}

func TestTaxRate(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	// New context. There should be no price.
	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

	terraPrice := sdk.NewDecWithPrec(166, 2)
	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(terraPrice) && err == nil)
}

func TestTaxCap(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	// New context. There should be no price.
	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

	terraPrice := sdk.NewDecWithPrec(166, 2)
	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
	require.True(t, tp.Equal(terraPrice) && err == nil)
}
