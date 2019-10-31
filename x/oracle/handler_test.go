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

	// Case 2: Normal MsgExchangeRatePrevote submission goes through
	salt := "1"
	bz, err := VoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)
	prevoteMsg := NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// // Case 3: Normal MsgExchangeRateVote submission goes through keeper.keeper.Addrs
	voteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 4: a non-validator sending an oracle message fails
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	bz, err = VoteHash(salt, randomExchangeRate, core.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	prevoteMsg = NewMsgExchangeRatePrevote("", core.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]))
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())
}

func TestPrevoteCheck(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	exchangeRatePrevoteMsg := NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, exchangeRatePrevoteMsg)
	require.True(t, res.IsOK())

	// Invalid exchange rate reveal period
	exchangeRateVoteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, exchangeRateVoteMsg)
	require.False(t, res.IsOK())

	input.Ctx = input.Ctx.WithBlockHeight(2)
	exchangeRateVoteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, exchangeRateVoteMsg)
	require.False(t, res.IsOK())

	// valid exchange rate reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	exchangeRateVoteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	res = h(input.Ctx, exchangeRateVoteMsg)
	require.True(t, res.IsOK())

}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.Nil(t, err)

	// Case 1: empty message
	bankMsg := MsgDelegateFeedConsent{}
	res := h(input.Ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 2.2: Normal Vote - without delegation
	voteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	// Case 2.3: Normal Vote - with delegation fails
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.False(t, res.IsOK())

	// Case 3: Normal MsgDelegateFeedConsent succeeds
	msg := NewMsgDelegateFeedConsent(keeper.ValAddrs[0], keeper.Addrs[1])
	res = h(input.Ctx, msg)
	require.True(t, res.IsOK())

	// Case 4.1: Normal Prevote - without delegation fails
	prevoteMsg = NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.False(t, res.IsOK())

	// Case 4.2: Normal Prevote - with delegation succeeds
	prevoteMsg = NewMsgExchangeRatePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())
	// Case 4.3: Normal Vote - without delegation fails
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.False(t, res.IsOK())

	// Case 4.4: Normal Vote - with delegation succeeds
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())
}
