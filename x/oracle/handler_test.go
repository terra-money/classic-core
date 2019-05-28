package oracle

import (
	"github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/terra-project/core/types/assets"

	"github.com/cosmos/cosmos-sdk/x/bank"
	cosmock "github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestOracleFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := bank.MsgSend{}
	res := h(input.ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal MsgPriceFeed submission goes through
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0], types.ValAddress(addrs[0]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 3: a non-validator sending an oracle message fails
	_, randoAddrs := cosmock.GeneratePrivKeyAddressPairs(1)
	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, randoAddrs[0], types.ValAddress(randoAddrs[0]))
	res = h(input.ctx, msg)
	require.False(t, res.IsOK())
}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	// Case 1: empty message
	bankMsg := MsgDelegateFeederPermission{}
	res := h(input.ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal Vote - without delegation
	priceMsg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0], types.ValAddress(addrs[0]))
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Vote - with delegation fails
	priceMsg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1], types.ValAddress(addrs[0]))
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateFeederPermission succeeds
	msg := NewMsgDelegateFeederPermission(types.ValAddress(addrs[0]), addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 4: Normal Vote - without delegation fails
	priceMsg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[2], types.ValAddress(addrs[0]))
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	// Case 5: Normal Vote - with delegation succeeds
	priceMsg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1], types.ValAddress(addrs[0]))
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())
}
