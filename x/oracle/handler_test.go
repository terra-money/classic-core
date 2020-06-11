package oracle

import (
	"fmt"
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
	_, err := h(input.Ctx, bankMsg)
	require.Error(t, err)

	// Case 2: Normal MsgExchangeRatePrevote submission goes through
	salt := "1"

	hash := GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// // Case 3: Normal MsgExchangeRateVote submission goes through keeper.keeper.Addrs
	voteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)

	// Case 4: a non-validator sending an oracle message fails
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	salt = "2"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, sdk.ValAddress(addrs[0]))

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]))
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)
}

func TestPrevoteVote(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	hash := GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	// Unintended denom prevote
	exchangeRatePrevoteMsg := NewMsgExchangeRatePrevote(hash, core.MicroCNYDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err := h(input.Ctx, exchangeRatePrevoteMsg)
	require.Error(t, err)

	// Unauthorized feeder
	exchangeRatePrevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRatePrevoteMsg)
	require.Error(t, err)

	// Valid prevote
	exchangeRatePrevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRatePrevoteMsg)
	require.NoError(t, err)

	// Invalid exchange rate reveal period
	exchangeRateVoteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRateVoteMsg)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(2)
	exchangeRateVoteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRateVoteMsg)
	require.Error(t, err)

	// Unauthorized feeder
	exchangeRateVoteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[1]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRateVoteMsg)
	require.Error(t, err)

	// Valid exchange rate reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	exchangeRateVoteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, exchangeRateVoteMsg)
	require.NoError(t, err)
}

func TestFeederDelegation(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	hash := GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	// Case 1: empty message
	bankMsg := MsgDelegateFeedConsent{}
	_, err := h(input.Ctx, bankMsg)
	require.Error(t, err)

	// Case 2: Normal Prevote - without delegation
	prevoteMsg := NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// Case 2.1: Normal Prevote - with delegation fails
	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)

	// Case 2.2: Normal Vote - without delegation
	voteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)

	// Case 2.3: Normal Vote - with delegation fails
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.Error(t, err)

	// Case 3: Normal MsgDelegateFeedConsent succeeds
	msg := NewMsgDelegateFeedConsent(keeper.ValAddrs[0], keeper.Addrs[1])
	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	// Case 4.1: Normal Prevote - without delegation fails
	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.Error(t, err)

	// Case 4.2: Normal Prevote - with delegation succeeds
	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// Case 4.3: Normal Vote - without delegation fails
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.Error(t, err)

	// Case 4.4: Normal Vote - with delegation succeeds
	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)
}

func TestAggregatePrevoteVote(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	exchangeRatesStr := fmt.Sprintf("1000.23%s,0.29%s,0.27%s", core.MicroKRWDenom, core.MicroUSDDenom, core.MicroSDRDenom)
	otherExchangeRateStr := fmt.Sprintf("1000.12%s,0.29%s,0.27%s", core.MicroKRWDenom, core.MicroUSDDenom, core.MicroUSDDenom)
	unintendedExchageRateStr := fmt.Sprintf("1000.23%s,0.29%s,0.27%s", core.MicroKRWDenom, core.MicroUSDDenom, core.MicroCNYDenom)
	invalidExchangeRateStr := fmt.Sprintf("1000.23%s,0.29%s,0.27", core.MicroKRWDenom, core.MicroUSDDenom)

	hash := GetAggregateVoteHash(salt, exchangeRatesStr, keeper.ValAddrs[0])

	aggregateExchangeRatePrevoteMsg := NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err := h(input.Ctx, aggregateExchangeRatePrevoteMsg)
	require.NoError(t, err)

	// Unauthorized feeder
	aggregateExchangeRatePrevoteMsg = NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRatePrevoteMsg)
	require.Error(t, err)

	// Invalid reveal period
	aggregateExchangeRateVoteMsg := NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Invalid reveal period
	input.Ctx = input.Ctx.WithBlockHeight(2)
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Other exchange rate with valid real period
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, otherExchangeRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Invalid exchange rate with valid real period
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, invalidExchangeRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Unauthorized feeder
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, invalidExchangeRateStr, sdk.AccAddress(keeper.Addrs[1]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Unintended denom vote
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, unintendedExchageRateStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.Error(t, err)

	// Valid exchange rate reveal submission
	input.Ctx = input.Ctx.WithBlockHeight(1)
	aggregateExchangeRateVoteMsg = NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, sdk.AccAddress(keeper.Addrs[0]), keeper.ValAddrs[0])
	_, err = h(input.Ctx, aggregateExchangeRateVoteMsg)
	require.NoError(t, err)
}
