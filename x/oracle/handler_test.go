package oracle

import (
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	cosmock "github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

var (
	randomPrice        = sdk.NewDecWithPrec(1049, 2) // swap rate
	anotherRandomPrice = sdk.NewDecWithPrec(4882, 2) // swap rate
)

func setup(t *testing.T) (testInput, sdk.Handler) {
	input := createTestInput(t)
	h := NewHandler(input.oracleKeeper)

	defaultOracleParams := DefaultParams()
	defaultOracleParams.VotePeriod = sdk.OneInt()
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	return input, h
}

func TestOracleFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := bank.MsgSend{}
	res := h(input.ctx, bankMsg)
	require.False(t, res.IsOK())

	// Case 2: Normal MsgPriceFeed submission goes through
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 3: a non-validator sending an oracle message fails
	_, randoAddrs := cosmock.GeneratePrivKeyAddressPairs(1)
	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, randoAddrs[0])
	res = h(input.ctx, msg)
	require.False(t, res.IsOK())
}

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	_, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[2])
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)

	// A new validator joins, we are now below threshold. Price update should now fail
	newValidator := mock.NewMockValidator(sdk.ValAddress(addrs[2].Bytes()), sdk.NewInt(30))
	input.valset.Validators = append(input.valset.Validators, newValidator)
	input.oracleKeeper.valset = input.valset

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[0])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[1])
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[2])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[0])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[2])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, price, anotherRandomPrice)
}

func TestOracleWhitelist(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroKRWDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	dropThreshold := input.oracleKeeper.GetParams(input.ctx).DropThreshold
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, randomPrice)

	msg := NewMsgPriceFeed(assets.MicroKRWDenom, randomPrice, addrs[0])
	h(input.ctx, msg)

	input.ctx = input.ctx.WithBlockHeight(1)
	for i := 0; i < int(dropThreshold.Int64())-1; i++ {
		EndBlocker(input.ctx, input.oracleKeeper)
	}

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroKRWDenom)
	require.Nil(t, err)
	require.Equal(t, price, randomPrice)

	// Going over dropthreshold should blacklist the price
	EndBlocker(input.ctx, input.oracleKeeper)

	_, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroKRWDenom)
	require.NotNil(t, err)
}
