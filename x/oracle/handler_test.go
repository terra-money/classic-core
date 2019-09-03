package oracle

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/keeper"
)

func TestOracleFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := bank.MsgSend{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal MsgPricePrevote submission goes through
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// // Case 3: Normal MsgPriceVote submission goes through keeper.keeper.Addrs
	voteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 4: a non-validator sending an oracle message fails
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	prevoteMsg = NewMsgPricePrevote("", core.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]))
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())
}

func TestPrevoteCheck(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	pricePrevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, pricePrevoteMsg)
	require.True(t, res.IsOK())

	// Invalid price reveal period
	priceVoteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, priceVoteMsg)
	require.False(t, res.IsOK())

	input.Ctx = input.Ctx.WithBlockHeight(2)
	priceVoteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, priceVoteMsg)
	require.False(t, res.IsOK())

	// valid price reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	priceVoteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, priceVoteMsg)
	require.True(t, res.IsOK())

}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	// Case 1: empty message
	bankMsg := MsgDelegateFeederPermission{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateFeederPermission succeeds
	msg := NewMsgDelegateFeederPermission(keeper.ValAddrs[0], keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.True(t, res.IsOK())

	// Case 4: Normal Prevote - without delegation fails
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 5: Normal Prevote - with delegation succeeds
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
}
