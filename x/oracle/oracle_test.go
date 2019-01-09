package oracle

import (
	"terra/types/assets"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/stake"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestOracle(t *testing.T) {
	mapp, keeper, sk, addrs, pubKeys, _ := getMockApp(t, 5)
	mapp.BeginBlock(abci.RequestBeginBlock{})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	//Set up validators
	stakeHandler := stake.NewHandler(sk)
	for i := 0; i < len(addrs); i++ {
		valCreateMsg := stake.NewMsgCreateValidator(
			sdk.ValAddress(addrs[i].Bytes()), pubKeys[i],
			sdk.NewInt64Coin(assets.LunaDenom, 10), stake.Description{},
			stake.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		)

		res := stakeHandler(ctx, valCreateMsg)
		require.True(t, res.IsOK())

		validator, _ := sk.GetValidator(ctx, sdk.ValAddress(addrs[i].Bytes()))
		validator.UpdateStatus(sk.GetPool(ctx), sdk.Bonded)
	}
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// Set up prices
	targetPrice := sdk.NewDecWithPrec(88, 2)
	observedPrice := sdk.NewDecWithPrec(77, 2)

	h := NewHandler(keeper)

	// Case 0: non-oracle message being sent fails
	msg := bank.MsgSend{}
	res := h(ctx, msg)
	require.False(t, res.IsOK())

	// Case 1: Normal pricefeedmsg submission goes through
	pfm := PriceFeedMsg{
		Denom:         assets.KRWDenom,
		TargetPrice:   targetPrice,
		ObservedPrice: observedPrice,
		Feeder:        addrs[0],
	}
	res = h(ctx, pfm)
	require.True(t, res.IsOK())

	// Case 1: a non-validator sending an oracle message fails
	_, randoAddrs := mock.GeneratePrivKeyAddressPairs(1)
	pfm.Feeder = randoAddrs[0]

	res = h(ctx, pfm)
	require.False(t, res.IsOK())

	// Case 2: sending a message for a non-whitelisted coin fails
	pfm.Feeder = addrs[0]
	pfm.Denom = "sketchyCoin"
	res = h(ctx, pfm)
	require.False(t, res.IsOK())

	// Case 3: less than the threshold has signs, msg fails
	pfm.Denom = assets.KRWDenom
	res = h(ctx, pfm)

	ctx = ctx.WithBlockHeight(1)

	require.Equal(t, sdk.ZeroDec(), keeper.GetPriceTarget(ctx, assets.KRWDenom))

	// Case 5: more than the threshold signs, msg succeeds
	//fmt.Printf("Block height : %d\n", ctx.BlockHeight())
	pfm.Feeder = addrs[1]
	res = h(ctx, pfm)
	pfm.Feeder = addrs[2]
	res = h(ctx, pfm)
	pfm.Feeder = addrs[3]
	res = h(ctx, pfm)
	pfm.Feeder = addrs[4]
	res = h(ctx, pfm)

	// votes := keeper.getVotes(ctx, assets.KRWDenom)
	// for _, vote := range votes {
	// 	fmt.Printf("%v \n", vote)
	// }

	// vals := sk.GetValidators(ctx, 10)
	// for _, v := range vals {
	// 	fmt.Printf("%v %v %v %v\n", v, v.GetPower(), v.GetTokens(), v.Status)
	// }

	ctx = ctx.WithBlockHeight(2)
	EndBlocker(ctx, keeper)

	//fmt.Printf("Block height : %d \n", ctx.BlockHeight())
	require.Equal(t, targetPrice, keeper.GetPriceTarget(ctx, assets.KRWDenom))
	require.Equal(t, observedPrice, keeper.GetPriceObserved(ctx, assets.KRWDenom))

	// Case 6: one of the previously signed validators are kicked out, msg should now fail

	// Case 7: more than the threshold signs, the msg now succeeds
}

// 	// cdc := makeCodec()

// 	// addr1 := []byte("addr1")
// 	// addr2 := []byte("addr2")
// 	// addr3 := []byte("addr3")
// 	// addr4 := []byte("addr4")
// 	// valset := &mock.ValidatorSet{[]mock.Validator{
// 	// 	{addr1, sdk.NewDec(7)},
// 	// 	{addr2, sdk.NewDec(7)},
// 	// 	{addr3, sdk.NewDec(1)},
// 	// }}

// 	// key := sdk.NewKVStoreKey("testkey")
// 	// ctx := defaultContext(key)

// 	// bz, err := json.Marshal(valset)
// 	// require.Nil(t, err)
// 	// ctx = ctx.WithBlockHeader(abci.Header{ValidatorsHash: bz})

// 	// tk := treasury.NewKeeper
// 	// ork := NewKeeper(key, cdc, valset, sdk.NewDecWithPrec(667, 3), 100) // 66.7%
// 	// // h := seqHandler(ork, key, sdk.CodespaceRoot)
// 	// h := NewHandler(oracle.Keeper)

// 	// // Nonmock.Validator signed, transaction failed
// 	// msg := Msg{seqOracle{0, 0}, []byte("randomguy")}
// 	// res := h(ctx, msg)
// 	// require.False(t, res.IsOK())
// 	// require.Equal(t, 0, getSequence(ctx, key))

// 	// // Less than 2/3 signed, msg not processed
// 	// msg.Signer = addr1
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 0, getSequence(ctx, key))

// 	// // Double signed, transaction failed
// 	// res = h(ctx, msg)
// 	// require.False(t, res.IsOK())
// 	// require.Equal(t, 0, getSequence(ctx, key))

// 	// // More than 2/3 signed, msg processed
// 	// msg.Signer = addr2
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // Already processed, transaction failed
// 	// msg.Signer = addr3
// 	// res = h(ctx, msg)
// 	// require.False(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // Less than 2/3 signed, msg not processed
// 	// msg = Msg{seqOracle{100, 1}, addr1}
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // More than 2/3 signed but payload is invalid
// 	// msg.Signer = addr2
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.NotEqual(t, "", res.Log)
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // Already processed, transaction failed
// 	// msg.Signer = addr3
// 	// res = h(ctx, msg)
// 	// require.False(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // Should handle mock.Validator set change
// 	// valset.AddValidator(mock.Validator{addr4, sdk.NewDec(12)})
// 	// bz, err = json.Marshal(valset)
// 	// require.Nil(t, err)
// 	// ctx = ctx.WithBlockHeader(abci.Header{ValidatorsHash: bz})

// 	// // Less than 2/3 signed, msg not processed
// 	// msg = Msg{seqOracle{1, 2}, addr1}
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // Less than 2/3 signed, msg not processed
// 	// msg.Signer = addr2
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 1, getSequence(ctx, key))

// 	// // More than 2/3 signed, msg processed
// 	// msg.Signer = addr4
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 2, getSequence(ctx, key))

// 	// // Should handle mock.Validator set change while oracle process is happening
// 	// msg = Msg{seqOracle{2, 3}, addr4}

// 	// // Less than 2/3 signed, msg not processed
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 2, getSequence(ctx, key))

// 	// // Signed mock.Validator is kicked out
// 	// valset.RemoveValidator(addr4)
// 	// bz, err = json.Marshal(valset)
// 	// require.Nil(t, err)
// 	// ctx = ctx.WithBlockHeader(abci.Header{ValidatorsHash: bz})

// 	// // Less than 2/3 signed, msg not processed
// 	// msg.Signer = addr1
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 2, getSequence(ctx, key))

// 	// // More than 2/3 signed, msg processed
// 	// msg.Signer = addr2
// 	// res = h(ctx, msg)
// 	// require.True(t, res.IsOK())
// 	// require.Equal(t, 3, getSequence(ctx, key))
// }

// func TestOracleRewards(t *testing.T) {

// }
