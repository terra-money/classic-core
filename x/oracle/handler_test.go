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

	// Case 2: Normal MsgPrevote submission goes through
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)
	prevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// // Case 3: Normal MsgVote submission goes through keeper.keeper.Addrs
	voteMsg := NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 4: a non-validator sending an oracle message fails
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	prevoteMsg = NewMsgPrevote("", core.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]))
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())
}

func TestPrevoteCheck(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	PrevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, PrevoteMsg)
	require.True(t, res.IsOK())

	// Invalid exchangeRate reveal period
	VoteMsg := NewMsgVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, VoteMsg)
	require.False(t, res.IsOK())

	input.Ctx = input.Ctx.WithBlockHeight(2)
	VoteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, VoteMsg)
	require.False(t, res.IsOK())

	// valid exchangeRate reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	VoteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, VoteMsg)
	require.True(t, res.IsOK())

}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	// Case 1: empty message
	bankMsg := MsgDelegateConsent{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 2.2: Normal Vote - without delegation
	voteMsg := NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 2.3: Normal Vote - with delegation fails
	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateConsent succeeds
	msg := NewMsgDelegateConsent(keeper.ValAddrs[0], keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.True(t, res.IsOK())

	// Case 4.1: Normal Prevote - without delegation fails
	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 4.2: Normal Prevote - with delegation succeeds
	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
	// Case 4.3: Normal Vote - without delegation fails
	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.False(t, res.IsOK())

	// Case 4.4: Normal Vote - with delegation succeeds
	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())
}
