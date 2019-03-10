package pay

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

// func TestGenesisIssuance(t *testing.T) {
// 	mapp, keeper, _, _, _ := getMockApp(t, 5)

// 	keeper.setIssuance(ctx, assets.SDRDenom, sdk.NewInt(10))

// }

func TestIssuance(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, 5)

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.setIssuance(ctx, assets.SDRDenom, sdk.NewInt(10))

	expectedIssuance := sdk.NewInt(10)
	keeper.GetIssuance(ctx, assets.SDRDenom, expectedIssuance)
	actualIssuance := keeper.GetIssuance(ctx, assets.SDRDenom, sdk.ZeroInt())

	require.Equal(t, expectedIssuance, actualIssuance, "Issuance does not match")
}

// func TestTaxRate(t *testing.T) {
// 	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
// 	mapp.BeginBlock(abci.RequestBeginBlock{})

// 	// New context. There should be no price.
// 	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
// 	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

// 	terraPrice := sdk.NewDecWithPrec(166, 2)
// 	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

// 	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
// 	require.True(t, tp.Equal(terraPrice) && err == nil)
// }

// func TestTaxCap(t *testing.T) {
// 	mapp, keeper, _, _, _, _ := getMockApp(t, 5)
// 	mapp.BeginBlock(abci.RequestBeginBlock{})
// 	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

// 	// New context. There should be no price.
// 	tp, err := keeper.GetPrice(ctx, assets.TerraDenom)
// 	require.True(t, tp.Equal(sdk.ZeroDec()) && err != nil)

// 	terraPrice := sdk.NewDecWithPrec(166, 2)
// 	keeper.setPrice(ctx, assets.TerraDenom, terraPrice)

// 	tp, err = keeper.GetPrice(ctx, assets.TerraDenom)
// 	require.True(t, tp.Equal(terraPrice) && err == nil)
// }
