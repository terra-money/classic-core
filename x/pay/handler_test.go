package pay

import (
	"terra/types/assets"
	"terra/types/util"
	"terra/x/market"
	"terra/x/mint"
	"terra/x/oracle"
	"terra/x/treasury"
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
	ctx            sdk.Context
	accKeeper      auth.AccountKeeper
	bankKeeper     bank.Keeper
	treasuryKeeper treasury.Keeper
	feeKeeper      auth.FeeCollectionKeeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyTreasury := sdk.NewKVStoreKey(treasury.StoreKey)
	keyFee := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFee, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)

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

	mintKeeper := mint.NewKeeper(cdc, keyMint, bankKeeper, accKeeper)
	var valset sdk.ValidatorSet
	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle,
		valset,
		paramsKeeper.Subspace(oracle.DefaultParamspace),
	)

	marketKeeper := market.NewKeeper(oracleKeeper, mintKeeper)

	treasuryKeeper := treasury.NewKeeper(
		cdc,
		keyTreasury,
		valset,
		mintKeeper,
		marketKeeper,
		paramsKeeper.Subspace(treasury.DefaultParamspace),
	)

	feeKeeper := auth.NewFeeCollectionKeeper(
		cdc, keyFee,
	)

	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(assets.SDRDenom, initAmt)})
		require.NoError(t, err)
	}

	return testInput{ctx, accKeeper, bankKeeper, treasuryKeeper, feeKeeper}
}

func TestHandlerMsgSendTransfersDisabled(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, false)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(5)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, amt)})

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	from := input.accKeeper.GetAccount(input.ctx, addrs[0])
	require.Equal(t, from.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, initAmt)})

	to := input.accKeeper.GetAccount(input.ctx, addrs[1])
	require.Equal(t, to.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, initAmt)})
}

func TestHandlerMsgSendTransfersEnabled(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)

	params := treasury.DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, params)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(5)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, amt)})

	res := handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	from := input.accKeeper.GetAccount(input.ctx, addrs[0])
	balance := initAmt.Sub(amt)
	require.Equal(t, from.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, balance)})

	to := input.accKeeper.GetAccount(input.ctx, addrs[1])
	balance = initAmt.Add(amt)
	require.Equal(t, to.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, balance)})
}

func TestHandlerMsgSendTax(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)
	params := treasury.DefaultParams()

	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.NewDecWithPrec(1, 3)) // 0.1%
	input.treasuryKeeper.SetParams(input.ctx, params)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(1000)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, amt)})

	handler(input.ctx, msg)

	taxCollected := input.feeKeeper.GetCollectedFees(input.ctx)
	taxRecorded := input.treasuryKeeper.PeekTaxProceeds(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(1))}, taxCollected)
	require.Equal(t, taxCollected, taxRecorded)

	remainingBalance := input.bankKeeper.GetCoins(input.ctx, addrs[0])
	require.Equal(t, remainingBalance, sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(4))}, "expected 4 SDR to be remaining")

	amt = sdk.NewInt(5)
	msg = bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, amt)})
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Clear coin balances
	input.bankKeeper.SetCoins(input.ctx, addrs[0], sdk.Coins{})
	input.bankKeeper.SetCoins(input.ctx, addrs[1], sdk.Coins{})

	// Give more coins
	input.bankKeeper.AddCoins(input.ctx, addrs[0], sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(5000))})

	// Reset tax cap
	params.TaxPolicy.Cap = sdk.NewInt64Coin(assets.SDRDenom, 2) // 2 SDR cap
	input.treasuryKeeper.SetParams(input.ctx, params)
	amt = sdk.NewInt(2000)
	msg = bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, amt)})
	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	remainingBalance = input.bankKeeper.GetCoins(input.ctx, addrs[0])
	expectedRemainingBalance := sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(2999))}
	receivedBalance := input.bankKeeper.GetCoins(input.ctx, addrs[1])
	expectedReceivedBalance := sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(2000))}
	taxCollected = input.feeKeeper.GetCollectedFees(input.ctx)
	expectedTaxCollected := sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(2))}

	require.Equal(t, expectedRemainingBalance, remainingBalance)
	require.Equal(t, expectedReceivedBalance, receivedBalance)
	require.Equal(t, expectedTaxCollected, taxCollected)
}

func TestHandlerMsgMultiSendTax(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)

	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.NewDecWithPrec(1, 2)) // 1%

	params := treasury.DefaultParams()
	params.TaxPolicy.Cap = sdk.NewInt64Coin(assets.SDRDenom, 2) // 2 SDR cap
	input.treasuryKeeper.SetParams(input.ctx, params)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)

	msg := bank.NewMsgMultiSend(
		[]bank.Input{
			bank.NewInput(addrs[0], sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(398))}),
			bank.NewInput(addrs[1], sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(14))}),
			bank.NewInput(addrs[2], sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(189))}),
		},
		[]bank.Output{
			bank.NewOutput(addrs[0], sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(601))}),
		},
	)

	res := handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	taxCollected := input.feeKeeper.GetCollectedFees(input.ctx)
	taxRecorded := input.treasuryKeeper.PeekTaxProceeds(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.Coins{sdk.NewCoin(assets.SDRDenom, sdk.NewInt(3))}, taxCollected)
	require.Equal(t, taxCollected, taxRecorded)

	acc1 := input.accKeeper.GetAccount(input.ctx, addrs[0])
	balance := initAmt.Add(sdk.NewInt(201))
	require.Equal(t, acc1.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, balance)})

	acc2 := input.accKeeper.GetAccount(input.ctx, addrs[1])
	balance = initAmt.Sub(sdk.NewInt(14))
	require.Equal(t, acc2.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, balance)})

	acc3 := input.accKeeper.GetAccount(input.ctx, addrs[2])
	balance = initAmt.Sub(sdk.NewInt(190))
	require.Equal(t, acc3.GetCoins(), sdk.Coins{sdk.NewCoin(assets.SDRDenom, balance)})
}
