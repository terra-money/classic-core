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

	// Case 2: Normal MsgPricePrevote submission goes through
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Case 3: Normal MsgPriceVote submission goes through
	voteMsg := NewMsgPriceVote(randomPrice, salt, assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]))
	res = h(input.ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 4: a non-validator sending an oracle message fails
	_, randoAddrs := cosmock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, types.ValAddress(randoAddrs[0]))
	require.Nil(t, err)

	prevoteMsg = NewMsgPricePrevote("", assets.MicroSDRDenom, randoAddrs[0], types.ValAddress(randoAddrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.False(t, res.IsOK())
}

func TestPrevoteCheck(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, types.ValAddress(addrs[0]))
	require.Nil(t, err)

	pricePrevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]))
	res := h(input.ctx, pricePrevoteMsg)
	require.True(t, res.IsOK())

	// Invalid price reveal period
	priceVoteMsg := NewMsgPriceVote(randomPrice, salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]))
	res = h(input.ctx, priceVoteMsg)
	require.False(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(2)
	priceVoteMsg = NewMsgPriceVote(randomPrice, salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]))
	res = h(input.ctx, priceVoteMsg)
	require.False(t, res.IsOK())

	// valid price reveal submission
	input.ctx = input.ctx.WithBlockHeight(1)
	priceVoteMsg = NewMsgPriceVote(randomPrice, salt, assets.MicroSDRDenom, sdk.AccAddress(addrs[0]), types.ValAddress(addrs[0]))
	res = h(input.ctx, priceVoteMsg)
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

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[0], types.ValAddress(addrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[1], types.ValAddress(addrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateFeederPermission succeeds
	msg := NewMsgDelegateFeederPermission(types.ValAddress(addrs[0]), addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 4: Normal Prevote - without delegation fails
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[2], types.ValAddress(addrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 5: Normal Prevote - with delegation succeeds
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), assets.MicroSDRDenom, addrs[1], types.ValAddress(addrs[0]))
	res = h(input.ctx, prevoteMsg)
	require.True(t, res.IsOK())
}
