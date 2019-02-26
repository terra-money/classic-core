package oracle

import (
	"terra/types/assets"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/mock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	trashDescription = stake.NewDescription("", "", "", "")
	trashCommission  = stake.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

func TestOracle(t *testing.T) {
	mapp, keeper, sk, addrs, pubKeys, _ := getMockApp(t, 6)
	mapp.BeginBlock(abci.RequestBeginBlock{})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	feedMsgs := []PriceFeedMsg{}

	//Set up validators
	stakeHandler := stake.NewHandler(sk)
	for i := 0; i < len(addrs)-1; i++ {
		valCreateMsg := stake.NewMsgCreateValidator(
			sdk.ValAddress(addrs[i].Bytes()), pubKeys[i],
			sdk.NewInt64Coin(assets.LunaDenom, 10), trashDescription, trashCommission,
		)

		res := stakeHandler(ctx, valCreateMsg)
		require.True(t, res.IsOK())

		validator, _ := sk.GetValidator(ctx, sdk.ValAddress(addrs[i].Bytes()))
		validator.UpdateStatus(sk.GetPool(ctx), sdk.Bonded)

		pfm := PriceFeedMsg{
			Denom:         assets.KRWDenom,
			TargetPrice:   sdk.NewDecWithPrec(98, 2),
			ObservedPrice: sdk.NewDecWithPrec(105, 2),
			Feeder:        addrs[i],
		}
		feedMsgs = append(feedMsgs, pfm)
	}
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	h := NewHandler(keeper)

	// Case 0: non-oracle message being sent fails
	msg := bank.MsgSend{}
	res := h(ctx, msg)
	require.False(t, res.IsOK())

	// Case 1: Normal pricefeedmsg submission goes through
	res = h(ctx, feedMsgs[0])
	require.True(t, res.IsOK())

	// Case 1: a non-validator sending an oracle message fails
	_, randoAddrs := mock.GeneratePrivKeyAddressPairs(1)
	feedMsgs[0].Feeder = randoAddrs[0]

	res = h(ctx, feedMsgs[0])
	require.False(t, res.IsOK())

	// Case 2: sending a message for a non-whitelisted coin fails
	feedMsgs[0].Feeder = addrs[0]
	feedMsgs[0].Denom = "sketchyCoin"
	res = h(ctx, feedMsgs[0])
	require.False(t, res.IsOK())

	// Case 3: less than the threshold signs, msg fails
	feedMsgs[0].Denom = assets.KRWDenom
	res = h(ctx, feedMsgs[0])

	ctx = ctx.WithBlockHeight(1)
	EndBlocker(ctx, keeper)

	require.Equal(t, sdk.ZeroDec(), keeper.GetPriceTarget(ctx, assets.KRWDenom))

	// Case 5: more than the threshold signs, msg succeeds
	h(ctx, feedMsgs[0])
	h(ctx, feedMsgs[1])
	h(ctx, feedMsgs[2])
	h(ctx, feedMsgs[3])

	ctx = ctx.WithBlockHeight(2)
	EndBlocker(ctx, keeper)

	require.Equal(t, sdk.NewDecWithPrec(98, 2), keeper.GetPriceTarget(ctx, assets.KRWDenom))
	require.Equal(t, sdk.NewDecWithPrec(105, 2), keeper.GetPriceObserved(ctx, assets.KRWDenom))

	// Case 6: A large validator joins, now below required majority. msg should now fail
	keeper.setPriceTarget(ctx, assets.KRWDenom, sdk.ZeroDec())
	keeper.setPriceObserved(ctx, assets.KRWDenom, sdk.ZeroDec())

	valAddr := sdk.ValAddress(addrs[5].Bytes())
	valKey := pubKeys[5]

	valCreateMsg := stake.NewMsgCreateValidator(valAddr, valKey, sdk.NewInt64Coin(assets.LunaDenom, 40), trashDescription, trashCommission)
	res = stakeHandler(ctx, valCreateMsg)
	require.True(t, res.IsOK())

	validator, _ := sk.GetValidator(ctx, valAddr)

	validator.UpdateStatus(sk.GetPool(ctx), sdk.Bonded)
	sk.ApplyAndReturnValidatorSetUpdates(ctx)

	h(ctx, feedMsgs[0])
	h(ctx, feedMsgs[1])
	h(ctx, feedMsgs[2])
	h(ctx, feedMsgs[3])

	ctx = ctx.WithBlockHeight(3)
	EndBlocker(ctx, keeper)

	require.Equal(t, sdk.ZeroDec(), keeper.GetPriceTarget(ctx, assets.KRWDenom))
	require.Equal(t, sdk.ZeroDec(), keeper.GetPriceObserved(ctx, assets.KRWDenom))
}
