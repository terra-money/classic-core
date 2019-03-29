package budget

import (
	"terra/types/assets"
	"terra/types/mock"
	"terra/types/util"
	"terra/x/mint"

	"github.com/cosmos/cosmos-sdk/x/staking"

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
)

var (
	pubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	addrs = []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
	}

	initAmt = sdk.NewInt(10000000)
	lunaAmt = sdk.NewInt(10)
)

type testInput struct {
	ctx          sdk.Context
	cdc          *codec.Codec
	mintKeeper   mint.Keeper
	bankKeeper   bank.Keeper
	budgetKeeper Keeper
	valset       mock.MockValset
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	bank.RegisterCodec(cdc)
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput() testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyBudget := sdk.NewKVStoreKey(StoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewKVStoreKey(staking.TStoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBudget, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeIAVL, db)

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err.Error())
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

	mintKeeper := mint.NewKeeper(
		cdc,
		keyMint,
		bankKeeper,
		accKeeper,
	)

	valset := mock.NewMockValSet()
	for _, addr := range addrs {

		err := mintKeeper.Mint(ctx, addr, sdk.NewCoin(assets.SDRDenom, initAmt))
		err2 := mintKeeper.Mint(ctx, addr, sdk.NewCoin(assets.LunaDenom, lunaAmt))

		if err != nil {
			panic(err.Error())
		}

		if err2 != nil {
			panic(err2.Error())
		}

		// Add validators
		validator := mock.NewMockValidator(sdk.ValAddress(addr.Bytes()), lunaAmt)
		valset.Validators = append(valset.Validators, validator)
	}

	budgetKeeper := NewKeeper(
		cdc, keyBudget, mintKeeper, valset,
		paramsKeeper.Subspace(DefaultParamspace),
	)

	InitGenesis(ctx, budgetKeeper, DefaultGenesisState())

	return testInput{ctx, cdc, mintKeeper, bankKeeper, budgetKeeper, valset}
}

func generateTestProgram(ctx sdk.Context, accounts ...sdk.AccAddress) Program {
	submitter := addrs[0]
	if len(accounts) > 0 {
		submitter = accounts[0]
	}

	executor := addrs[1]
	if len(accounts) > 1 {
		executor = accounts[1]
	}

	return NewProgram("testTitle", "testDescription", submitter, executor, util.GetEpoch(ctx).Int64())
}
