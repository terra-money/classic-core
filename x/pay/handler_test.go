package pay

import (
	"fmt"
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/mint"
	"github.com/terra-project/core/x/oracle"
	"github.com/terra-project/core/x/treasury"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	valConsPubKeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}

	valConsAddrs = []sdk.ConsAddress{
		sdk.ConsAddress(valConsPubKeys[0].Address()),
		sdk.ConsAddress(valConsPubKeys[1].Address()),
		sdk.ConsAddress(valConsPubKeys[2].Address()),
	}

	uSDRAmount  = sdk.NewInt(1005).MulRaw(assets.MicroUnit)
	uLunaAmount = sdk.NewInt(1005).MulRaw(assets.MicroUnit)
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
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	tKeyDistr := sdk.NewTransientStoreKey(distr.TStoreKey)
	keyFeeCollection := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyMarket := sdk.NewKVStoreKey(market.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyDistr, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyFeeCollection, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)

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

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking, tKeyStaking,
		bankKeeper, paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	feeCollectionKeeper := auth.NewFeeCollectionKeeper(
		cdc,
		keyFeeCollection,
	)

	distrKeeper := distr.NewKeeper(
		cdc, keyDistr, paramsKeeper.Subspace(distr.DefaultParamspace),
		bankKeeper, &stakingKeeper, feeCollectionKeeper, distr.DefaultCodespace,
	)

	stakingKeeper.SetPool(ctx, staking.InitialPool())
	stakingKeeper.SetParams(ctx, staking.DefaultParams())

	mintKeeper := mint.NewKeeper(
		cdc,
		keyMint,
		stakingKeeper,
		bankKeeper,
		accKeeper,
	)

	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle,
		mintKeeper,
		distrKeeper,
		feeCollectionKeeper,
		&stakingKeeper,
		paramsKeeper.Subspace(oracle.DefaultParamspace),
	)

	marketKeeper := market.NewKeeper(cdc, keyMarket, oracleKeeper, mintKeeper,
		paramsKeeper.Subspace(market.DefaultParamspace))
	marketKeeper.SetParams(ctx, market.DefaultParams())

	treasuryKeeper := treasury.NewKeeper(
		cdc,
		keyTreasury,
		stakingKeeper.GetValidatorSet(),
		mintKeeper,
		marketKeeper,
		paramsKeeper.Subspace(treasury.DefaultParamspace),
	)

	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, uSDRAmount)})
		require.NoError(t, err)
	}

	return testInput{ctx, accKeeper, bankKeeper, treasuryKeeper, feeCollectionKeeper}
}

func TestHandlerMsgSendTransfersDisabled(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, false)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(5)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, amt)})

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	from := input.accKeeper.GetAccount(input.ctx, addrs[0])
	require.Equal(t, from.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, uSDRAmount)})

	to := input.accKeeper.GetAccount(input.ctx, addrs[1])
	require.Equal(t, to.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, uSDRAmount)})
}

func TestHandlerMsgSendTransfersEnabled(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)

	params := treasury.DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, params)
	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.ZeroDec()) // 0.0%

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(5).MulRaw(assets.MicroUnit)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, amt)})

	res := handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	from := input.accKeeper.GetAccount(input.ctx, addrs[0])
	balance := uSDRAmount.Sub(amt)
	require.Equal(t, from.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, balance)})

	to := input.accKeeper.GetAccount(input.ctx, addrs[1])
	balance = uSDRAmount.Add(amt)
	require.Equal(t, to.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, balance)})
}

func TestHandlerMsgSendTax(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)
	params := treasury.DefaultParams()

	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.NewDecWithPrec(1, 3)) // 0.1%
	input.treasuryKeeper.SetParams(input.ctx, params)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(1000).MulRaw(assets.MicroUnit)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, amt)})

	handler(input.ctx, msg)

	taxCollected := input.feeKeeper.GetCollectedFees(input.ctx)
	taxRecorded := input.treasuryKeeper.PeekTaxProceeds(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1).MulRaw(assets.MicroUnit))}, taxCollected)
	require.Equal(t, taxCollected, taxRecorded)

	remainingBalance := input.bankKeeper.GetCoins(input.ctx, addrs[0])
	require.Equal(t, remainingBalance, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))}, "expected 4 SDR to be remaining")

	amt = sdk.NewInt(5).MulRaw(assets.MicroUnit)
	msg = bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, amt)})
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Clear coin balances
	err := input.bankKeeper.SetCoins(input.ctx, addrs[0], sdk.Coins{})
	require.Nil(t, err)
	err = input.bankKeeper.SetCoins(input.ctx, addrs[1], sdk.Coins{})
	require.Nil(t, err)

	// Give more coins
	_, _, err = input.bankKeeper.AddCoins(input.ctx, addrs[0], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(5000).MulRaw(assets.MicroUnit))})
	require.Nil(t, err)

	// Reset tax cap
	params.TaxPolicy.Cap = sdk.NewInt64Coin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit).Int64()) // 2 SDR cap
	input.treasuryKeeper.SetParams(input.ctx, params)
	amt = sdk.NewInt(2000).MulRaw(assets.MicroUnit)
	msg = bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, amt)})
	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	remainingBalance = input.bankKeeper.GetCoins(input.ctx, addrs[0])
	expectedRemainingBalance := sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(2999).MulRaw(assets.MicroUnit))}
	receivedBalance := input.bankKeeper.GetCoins(input.ctx, addrs[1])
	expectedReceivedBalance := sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(2000).MulRaw(assets.MicroUnit))}
	taxCollected = input.feeKeeper.GetCollectedFees(input.ctx)
	expectedTaxCollected := sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewDecFromIntWithPrec(sdk.NewInt(2005000), 6).MulInt64(assets.MicroUnit).TruncateInt())}

	fmt.Println(taxCollected)
	require.Equal(t, expectedRemainingBalance, remainingBalance)
	require.Equal(t, expectedReceivedBalance, receivedBalance)
	require.Equal(t, expectedTaxCollected, taxCollected)
}

func TestLunaSendTax(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)
	params := treasury.DefaultParams()

	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.NewDecWithPrec(1, 3)) // 0.1%
	input.treasuryKeeper.SetParams(input.ctx, params)

	_, _, err := input.bankKeeper.AddCoins(input.ctx, addrs[0], sdk.Coins{sdk.NewCoin(assets.MicroLunaDenom, uLunaAmount)})
	require.NoError(t, err)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)
	amt := sdk.NewInt(1000).MulRaw(assets.MicroUnit)
	msg := bank.NewMsgSend(addrs[0], addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroLunaDenom, amt)})

	handler(input.ctx, msg)

	taxCollected := input.feeKeeper.GetCollectedFees(input.ctx)
	taxRecorded := input.treasuryKeeper.PeekTaxProceeds(input.ctx, util.GetEpoch(input.ctx))
	require.Equal(t, sdk.Coins{}, taxCollected)
	require.Equal(t, taxCollected, taxRecorded)

	require.Equal(t, amt, input.bankKeeper.GetCoins(input.ctx, addrs[1]).AmountOf(assets.MicroLunaDenom))
	require.Equal(t, uLunaAmount.Sub(amt), input.bankKeeper.GetCoins(input.ctx, addrs[0]).AmountOf(assets.MicroLunaDenom))
}

func TestHandlerMsgMultiSendTax(t *testing.T) {
	input := createTestInput(t)
	input.bankKeeper.SetSendEnabled(input.ctx, true)

	input.treasuryKeeper.SetTaxRate(input.ctx, sdk.NewDecWithPrec(1, 2)) // 1%

	params := treasury.DefaultParams()
	params.TaxPolicy.Cap = sdk.NewInt64Coin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit).Int64()) // 2 SDR cap
	input.treasuryKeeper.SetParams(input.ctx, params)

	handler := NewHandler(input.bankKeeper, input.treasuryKeeper, input.feeKeeper)

	msg := bank.NewMsgMultiSend(
		[]bank.Input{
			bank.NewInput(addrs[0], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(398).MulRaw(assets.MicroUnit))}),
			bank.NewInput(addrs[1], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(14).MulRaw(assets.MicroUnit))}),
			bank.NewInput(addrs[2], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(189).MulRaw(assets.MicroUnit))}),
		},
		[]bank.Output{
			bank.NewOutput(addrs[0], sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(601).MulRaw(assets.MicroUnit))}),
		},
	)

	res := handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	taxCollected := input.feeKeeper.GetCollectedFees(input.ctx)
	taxRecorded := input.treasuryKeeper.PeekTaxProceeds(input.ctx, util.GetEpoch(input.ctx))

	require.Equal(t, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, sdk.NewDecFromIntWithPrec(sdk.NewInt(403), 2).MulInt64(assets.MicroUnit).TruncateInt())}, taxCollected)
	require.Equal(t, taxCollected, taxRecorded)

	acc1 := input.accKeeper.GetAccount(input.ctx, addrs[0])
	balance := uSDRAmount.Add(sdk.NewInt(201).MulRaw(assets.MicroUnit))
	require.Equal(t, acc1.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, balance)})

	acc2 := input.accKeeper.GetAccount(input.ctx, addrs[1])
	balance = uSDRAmount.Sub(sdk.NewDecFromIntWithPrec(sdk.NewInt(1414), 2).MulInt64(assets.MicroUnit).TruncateInt())
	require.Equal(t, acc2.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, balance)})

	acc3 := input.accKeeper.GetAccount(input.ctx, addrs[2])
	balance = uSDRAmount.Sub(sdk.NewDecFromIntWithPrec(sdk.NewInt(19089), 2).MulInt64(assets.MicroUnit).TruncateInt())
	require.Equal(t, acc3.GetCoins(), sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, balance)})
}
