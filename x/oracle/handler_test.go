package oracle

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

var (
	randomPrice        = sdk.NewDecWithPrec(1049, 2)
	anotherRandomPrice = sdk.NewDecWithPrec(4882, 2)
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
	msg := NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[0])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Case 3: a non-validator sending an oracle message fails
	_, randoAddrs := mock.GeneratePrivKeyAddressPairs(1)
	msg = NewMsgPriceFeed(assets.SDRDenom, randomPrice, randoAddrs[0])
	res = h(input.ctx, msg)
	require.False(t, res.IsOK())
}

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	_, err := input.oracleKeeper.GetPrice(input.ctx, assets.SDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	msg = NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[0])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[1])
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetPrice(input.ctx, assets.SDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)

	// A new validator joins, we are now below threshold. Price update should now fail
	newValidator := NewMockValidator(sdk.ValAddress(addrs[2].Bytes()), sdk.NewInt(30))
	input.valset.validators = append(input.valset.validators, newValidator)
	input.oracleKeeper.valset = input.valset

	msg = NewMsgPriceFeed(assets.SDRDenom, anotherRandomPrice, addrs[0])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.SDRDenom, anotherRandomPrice, addrs[1])
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err = input.oracleKeeper.GetPrice(input.ctx, assets.SDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.SDRDenom, randomPrice, addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.SDRDenom, anotherRandomPrice, addrs[0])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.SDRDenom, anotherRandomPrice, addrs[1])
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetPrice(input.ctx, assets.SDRDenom)
	require.Nil(t, err)
	require.Equal(t, price, anotherRandomPrice)
}

func TestOracleWhitelist(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.KRWDenom, randomPrice, addrs[0])
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	dropThreshold := input.oracleKeeper.GetParams(input.ctx).DropThreshold
	input.oracleKeeper.SetPrice(input.ctx, assets.KRWDenom, randomPrice)

	msg := NewMsgPriceFeed(assets.KRWDenom, randomPrice, addrs[0])
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.KRWDenom, randomPrice, addrs[1])
	h(input.ctx, msg)

	input.ctx = input.ctx.WithBlockHeight(1)
	for i := 0; i < int(dropThreshold.Int64())-1; i++ {
		EndBlocker(input.ctx, input.oracleKeeper)
	}

	price, err := input.oracleKeeper.GetPrice(input.ctx, assets.KRWDenom)
	require.Nil(t, err)
	require.Equal(t, price, randomPrice)

	// Going over dropthreshold should blacklist the price
	for i := 0; i < int(dropThreshold.Int64())+1; i++ {
		EndBlocker(input.ctx, input.oracleKeeper)
	}

	price, err = input.oracleKeeper.GetPrice(input.ctx, assets.KRWDenom)
	require.NotNil(t, err)
}
