package market

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/keeper"
)

func TestMarketFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := bank.MsgSend{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal MsgSwap submission goes through
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	prevoteMsg := NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
}

func TestSwapMsg(t *testing.T) {
	input, h := setup(t)

	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	prevoteMsg := NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	res := h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
}
