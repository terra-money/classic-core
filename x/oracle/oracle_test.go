package oracle

import (
	"terra/types/assets"
	"terra/types/tax"
	"terra/x/treasury"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	//RegisterCodec(mapp.Cdc)

	keyGlobalParams := sdk.NewKVStoreKey("params")
	tkeyGlobalParams := sdk.NewTransientStoreKey("transient_params")
	keyStake := sdk.NewKVStoreKey("stake")
	tkeyStake := sdk.NewTransientStoreKey("transient_stake")
	keyTreasury := sdk.NewKVStoreKey("treasury")
	keyOracle := sdk.NewKVStoreKey("oracle")
	keyBank := sdk.NewKVStoreKey("bank")
	keyFeeCollection := sdk.NewKVStoreKey("fee")

	pk := params.NewKeeper(mapp.Cdc, keyGlobalParams, tkeyGlobalParams)
	fck := auth.NewFeeCollectionKeeper(mapp.Cdc, keyFeeCollection)
	ck := tax.NewBaseKeeper(keyBank, mapp.Cdc, mapp.AccountKeeper, fck)
	sk := stake.NewKeeper(mapp.Cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
	tk := treasury.NewKeeper(keyTreasury, mapp.Cdc, ck)
	keeper := NewKeeper(keyOracle, mapp.Cdc, tk, sk.GetValidatorSet(), pk.Subspace(DefaultParamspace))

	mapp.Router().AddRoute("oracle", NewHandler(keeper))
	//mapp.QueryRouter().AddRoute("oracle", NewQuerier(keeper))

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(assets.LunaDenom, 42)})

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk, addrs, pubKeys))

	require.NoError(t, mapp.CompleteSetup(keyStake, tkeyStake,
		keyTreasury,
		keyOracle, keyBank, keyFeeCollection,
		keyGlobalParams, tkeyGlobalParams))

	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}

// gov and stake endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// oracle and stake initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakeKeeper stake.Keeper, addrs []sdk.AccAddress, pubKeys []crypto.PubKey) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		defaultDescription := stake.Description{
			Moniker:  "",
			Identity: "",
			Website:  "",
			Details:  "",
		}

		stakeGenesis := stake.GenesisState{
			Pool: stake.InitialPool(),
			Params: stake.Params{
				UnbondingTime: 60 * 60 * 24 * 3 * time.Second,
				MaxValidators: 100,
				BondDenom:     assets.LunaDenom,
			},
			Validators: []stake.Validator{
				stake.NewValidator(sdk.ValAddress(addrs[0].Bytes()), pubKeys[0], defaultDescription),
				stake.NewValidator(sdk.ValAddress(addrs[1].Bytes()), pubKeys[1], defaultDescription),
				stake.NewValidator(sdk.ValAddress(addrs[2].Bytes()), pubKeys[2], defaultDescription),
			},
		}
		stakeGenesis.Pool.LooseTokens = sdk.NewDec(100000)

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		InitGenesis(ctx, keeper, DefaultGenesisState())
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

func defaultContext(keys ...sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, nil)
	return ctx
}

// func TestOracle(t *testing.T) {
// 	mapp, keeper, _, addrs, _, _ := getMockApp(t, 5)
// 	ctx := mapp.Contex

// 	//valset := sk.GetValidatorSet()

// 	h := NewHandler(keeper)

// 	// Case 0: non-oracle message being sent fails
// 	msg := bank.MsgSend{}
// 	res := h(mapp.Con, msg)
// 	require.False(t, res.IsOK())

// 	// Case 1: Normal pricefeedmsg submission goes through
// 	pfm := PriceFeedMsg{
// 		Denom:         assets.KRWDenom,
// 		TargetPrice:   sdk.OneDec(),
// 		ObservedPrice: sdk.OneDec(),
// 		Feeder:        addrs[0],
// 	}
// 	res = h(ctx, pfm)
// 	require.True(t, res.IsOK())

// 	// Case 1: a non-validator sending an oracle message fails

// 	// msg := Msg{seqOracle{0, 0}, []byte("randomguy")}
// 	// res := h(ctx, msg)

// 	// Case 2: sending a message for a non-whitelisted coin fails

// 	// Case 3: less than the threshold has signs, msg fails

// 	// Case 4: double signing, msg fails

// 	// Case 5: more than the threshold signs, msg succeeds

// 	// Case 6: one of the previously signed validators are kicked out, msg should now fail

// 	// Case 7: more than the threshold signs, the msg now succeeds

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
