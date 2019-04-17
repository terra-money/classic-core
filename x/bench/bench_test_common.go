package bench

import (
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	"github.com/terra-project/core/x/budget"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/mint"
	"github.com/terra-project/core/x/oracle"
	"github.com/terra-project/core/x/treasury"

	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const numOfValidators = 100

var (
	pubKeys [numOfValidators]crypto.PubKey
	addrs   [numOfValidators]sdk.AccAddress

	mLunaAmt = sdk.NewInt(10000000000).MulRaw(assets.MicroUnit)
	mSDRAmt  = sdk.NewInt(10000000000).MulRaw(assets.MicroUnit)
)

type testInput struct {
	ctx            sdk.Context
	cdc            *codec.Codec
	bankKeeper     bank.Keeper
	budgetKeeper   budget.Keeper
	oracleKeeper   oracle.Keeper
	marketKeeper   market.Keeper
	mintKeeper     mint.Keeper
	treasuryKeeper treasury.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	bank.RegisterCodec(cdc)
	treasury.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// createTestInput common test code for bench test
func createTestInput() testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyBudget := sdk.NewKVStoreKey(budget.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyTreasury := sdk.NewKVStoreKey(treasury.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBudget, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeTransient, db)

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

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

	stakingKeeper.SetPool(ctx, staking.InitialPool())
	stakingKeeper.SetParams(ctx, staking.DefaultParams())

	mintKeeper := mint.NewKeeper(
		cdc,
		keyMint,
		stakingKeeper,
		bankKeeper,
		accKeeper,
	)

	valset := mock.NewMockValSet()
	for i := 0; i < 100; i++ {
		pubKeys[i] = secp256k1.GenPrivKey().PubKey()
		addrs[i] = sdk.AccAddress(pubKeys[i].Address())

		err := mintKeeper.Mint(ctx, addrs[i], sdk.NewCoin(assets.MicroLunaDenom, mLunaAmt))
		if err != nil {
			panic(err)
		}

		err = mintKeeper.Mint(ctx, addrs[i], sdk.NewCoin(assets.MicroSDRDenom, mSDRAmt))
		if err != nil {
			panic(err)
		}

		// Add validators
		validator := mock.NewMockValidator(sdk.ValAddress(addrs[i].Bytes()), mLunaAmt)
		valset.Validators = append(valset.Validators, validator)
	}

	budgetKeeper := budget.NewKeeper(
		cdc, keyBudget, mintKeeper, valset,
		paramsKeeper.Subspace(budget.DefaultParamspace),
	)

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

	budget.InitGenesis(ctx, budgetKeeper, budget.DefaultGenesisState())
	oracle.InitGenesis(ctx, oracleKeeper, oracle.DefaultGenesisState())
	treasury.InitGenesis(ctx, treasuryKeeper, treasury.DefaultGenesisState())

	return testInput{ctx, cdc, bankKeeper, budgetKeeper, oracleKeeper, marketKeeper, mintKeeper, treasuryKeeper}
}
