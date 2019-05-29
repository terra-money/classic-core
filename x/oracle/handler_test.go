package oracle

import (
	"encoding/hex"
	"testing"

	"github.com/terra-project/core/types/assets"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	salt := "1"
	bz, err := VoteHash("1", randomPrice, assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)
	msg := NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 3: a non-validator sending an oracle message fails
	_, randoAddrs := cosmock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	bz, err = VoteHash("1", randomPrice, assets.MicroSDRDenom, types.ValAddress(randoAddrs[0]))
	require.Nil(t, err)
	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, randoAddrs[0], types.ValAddress(randoAddrs[0]), randomPrice)
	res = h(input.ctx, msg)
	require.False(t, res.IsOK())
}

func TestPrevoteCheck(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)

	priceMsg := NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]), sdk.ZeroDec())
	res := h(input.ctx, priceMsg)
	require.True(t, res.IsOK())

	// Invalid price reveal period
	priceMsg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(2)
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	// valid price reveal submission
	input.ctx = input.ctx.WithBlockHeight(1)
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())

	// valid hash change submission with zero price
	input.ctx = input.ctx.WithBlockHeight(1)
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]), sdk.ZeroDec())
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())

}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)

	// Case 1: empty message
	bankMsg := MsgDelegateFeederPermission{}
	res := h(input.ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal Vote - without delegation
	priceMsg := NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Vote - with delegation fails
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, addrs[1], types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateFeederPermission succeeds
	msg := NewMsgDelegateFeederPermission(types.ValAddress(addrs[0]), addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 4: Normal Vote - without delegation fails
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, addrs[2], types.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, priceMsg)
	require.False(t, res.IsOK())

	// Case 5: Normal Vote - with delegation succeeds
	priceMsg = NewMsgPriceFeed(hex.EncodeToString(bz), salt, assets.MicroSDRDenom, addrs[1], types.ValAddress(addrs[0]), sdk.ZeroDec())
	res = h(input.ctx, priceMsg)
	require.True(t, res.IsOK())
}
