package market

import (
	"terra/types/assets"
	"terra/x/mint"
	"terra/x/oracle"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	initAmt = sdk.NewInt(1005)
)

type testInput struct {
	ctx          sdk.Context
	accKeeper    auth.AccountKeeper
	bankKeeper   bank.Keeper
	marketKeeper Keeper
	oracleKeeper oracle.Keeper
	mintKeeper   mint.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	bank.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyMarket := sdk.NewKVStoreKey(StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accKeeper := auth.NewAccountKeeper(
		cdc,
		keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	bankKeeper := bank.NewBaseKeeper(
		accKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	var valset sdk.ValidatorSet
	oracleKeeper := oracle.NewKeeper(
		cdc, keyOracle, valset,
		paramsKeeper.Subspace(oracle.DefaultParamspace),
	)

	mintKeeper := mint.NewKeeper(
		cdc,
		keyMint,
		bankKeeper,
		accKeeper,
	)

	marketKeeper := Keeper{
		oracleKeeper, mintKeeper,
	}

	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(assets.SDRDenom, initAmt)})
		require.NoError(t, err)
	}

	return testInput{ctx, accKeeper, bankKeeper, marketKeeper, oracleKeeper, mintKeeper}
}

func TestHandlerMsgSwapValidPrice(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerAmt := sdk.NewInt(2)
	offerCoin := sdk.NewCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(4)
	askLunaPrice := sdk.NewDec(8)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)

	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	swapAmount := offerLunaPrice.Quo(askLunaPrice).MulInt(offerAmt).TruncateInt()
	trader := input.accKeeper.GetAccount(input.ctx, addrs[0])
	require.Equal(t, trader.GetCoins().AmountOf(offerCoin.Denom), initAmt.Sub(offerAmt))
	require.Equal(t, trader.GetCoins().AmountOf(askDenom), swapAmount)
}

func TestHandlerMsgSwapNoBalance(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	// Try to swap a coin I don't have at all
	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.CNYDenom, sdk.OneInt()), assets.GBPDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Try to swap a coin I don't have enough of
	msg.OfferCoin = sdk.NewCoin(assets.SDRDenom, initAmt.Add(sdk.OneInt()))
	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapRecursion(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.SDRDenom, sdk.OneInt()), assets.SDRDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapTooSmall(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerAmt := sdk.NewInt(100)
	offerCoin := sdk.NewCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(1)
	askLunaPrice := sdk.NewDecWithPrec(1001, 1)

	// Set oracle price
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Reset oracle price
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	askLunaPrice = sdk.NewDecWithPrec(1000, 1)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)
}
