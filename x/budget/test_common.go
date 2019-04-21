package budget

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/mint"

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

	uSDRAmt  = sdk.NewInt(1005 * assets.MicroUnit)
	uLunaAmt = sdk.NewInt(10 * assets.MicroUnit)
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

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyBudget := sdk.NewKVStoreKey(StoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)

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
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeTransient, db)

	if err := ms.LoadLatestVersion(); err != nil {
		require.Nil(t, err)
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
	for _, addr := range addrs {
		err := mintKeeper.Mint(ctx, addr, sdk.NewCoin(assets.MicroSDRDenom, uSDRAmt))
		err2 := mintKeeper.Mint(ctx, addr, sdk.NewCoin(assets.MicroLunaDenom, uLunaAmt))

		if err != nil {
			require.Nil(t, err)
		}

		if err2 != nil {
			require.Nil(t, err2)
		}

		// Add validators
		validator := mock.NewMockValidator(sdk.ValAddress(addr.Bytes()), uLunaAmt)
		valset.Validators = append(valset.Validators, validator)
	}

	budgetKeeper := NewKeeper(
		cdc, keyBudget, mintKeeper, valset,
		paramsKeeper.Subspace(DefaultParamspace),
	)

	InitGenesis(ctx, budgetKeeper, DefaultGenesisState())

	return testInput{ctx, cdc, mintKeeper, bankKeeper, budgetKeeper, valset}
}

func generateTestProgram(ctx sdk.Context, budgetKeeper Keeper, accounts ...sdk.AccAddress) Program {
	submitter := addrs[0]
	if len(accounts) > 0 {
		submitter = accounts[0]
	}

	executor := addrs[1]
	if len(accounts) > 1 {
		executor = accounts[1]
	}

	testProgramID := budgetKeeper.NewProgramID(ctx)

	return NewProgram(testProgramID, "testTitle", "testDescription", submitter, executor, util.GetEpoch(ctx).Int64())
}
