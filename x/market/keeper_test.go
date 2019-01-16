package market

import (
	"fmt"
	"terra/types/assets"
	"terra/x/oracle"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestSwap(t *testing.T) {

	// Addr 0 is the only validator
	mapp, keeper, ok, addrs, _, _ := getMockApp(t, 5)
	oh := oracle.NewHandler(ok)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	// No price set to the oracle. SwapCoins should fail.
	_, err := keeper.SwapCoins(ctx, sdk.NewInt64Coin(assets.KRWDenom, 10), assets.USDDenom)
	require.NotNil(t, err)

	// Added exchange rate for offer denom. Test should still fail.
	pfm := oracle.NewPriceFeedMsg(assets.KRWDenom, sdk.OneDec(), sdk.OneDec(), addrs[0])
	res := oh(ctx, pfm)
	require.True(t, res.IsOK())
	ctx = ctx.WithBlockHeight(1000000)
	oracle.EndBlocker(ctx, ok)

	_, err = keeper.SwapCoins(ctx, sdk.NewInt64Coin(assets.KRWDenom, 10), assets.USDDenom)
	require.NotNil(t, err)

	// Now ask denom should have exchange rate set. Test should still fail
	_, err = keeper.SwapCoins(ctx, sdk.NewInt64Coin(assets.USDDenom, 10), assets.KRWDenom)
	require.NotNil(t, err)

	// Both denoms set. should succeed.
	pfm = oracle.NewPriceFeedMsg(assets.USDDenom, sdk.OneDec(), sdk.OneDec(), addrs[0])
	res = oh(ctx, pfm)
	require.True(t, res.IsOK())

	pfm = oracle.NewPriceFeedMsg(assets.KRWDenom, sdk.OneDec(), sdk.OneDec(), addrs[0])
	res = oh(ctx, pfm)
	require.True(t, res.IsOK())

	ctx = ctx.WithBlockHeight(2000000)
	oracle.EndBlocker(ctx, ok)

	fmt.Printf("%v, %v\n", ok.GetPriceTarget(ctx, assets.KRWDenom), ok.GetPriceTarget(ctx, assets.USDDenom))

	_, err = keeper.SwapCoins(ctx, sdk.NewInt64Coin(assets.USDDenom, 10), assets.KRWDenom)
	require.Nil(t, err)
}
