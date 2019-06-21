package plan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/mint"
	"github.com/terra-project/core/x/oracle"

	"time"

	abci "github.com/tendermint/tendermint/abci/types"
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
	oracleDecPrecision = 8
)

func setup(t *testing.T) testInput {
	input := createTestInput(t)

	defaultOracleParams := oracle.DefaultParams()
	defaultOracleParams.VotePeriod = int64(1) // Set to one block for convinience
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	defaultMarketParams := market.DefaultParams()
	defaultMarketParams.DailyLunaDeltaCap = sdk.NewDecWithPrec(5, 3)
	input.marketKeeper.SetParams(input.ctx, defaultMarketParams)

	return input
}

type testInput struct {
	ctx          sdk.Context
	cdc          *codec.Codec
	accKeeper    auth.AccountKeeper
	oracleKeeper oracle.Keeper
	marketKeeper market.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	types.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyMarket := sdk.NewKVStoreKey(market.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tKeyStaking := sdk.NewKVStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	tKeyDistr := sdk.NewTransientStoreKey(distr.TStoreKey)
	keyFeeCollection := sdk.NewKVStoreKey(auth.FeeStoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeIAVL, db)
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

	mintKeeper := mint.NewKeeper(
		cdc,
		keyMint,
		stakingKeeper,
		bankKeeper,
		accKeeper,
	)

	stakingKeeper.SetPool(ctx, staking.InitialPool())
	stakingParams := staking.DefaultParams()
	stakingParams.BondDenom = assets.MicroLunaDenom
	stakingKeeper.SetParams(ctx, stakingParams)

	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle,
		mintKeeper,
		distrKeeper,
		feeCollectionKeeper,
		stakingKeeper.GetValidatorSet(),
		paramsKeeper.Subspace(oracle.DefaultParamspace),
	)

	marketKeeper := market.NewKeeper(
		cdc,
		keyMarket,
		oracleKeeper,
		mintKeeper,
		paramsKeeper.Subspace(market.DefaultParamspace),
	)

	return testInput{ctx, cdc, accKeeper, oracleKeeper, marketKeeper}
}
