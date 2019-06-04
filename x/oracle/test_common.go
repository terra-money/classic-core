package oracle

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"

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

	randomPrice        = sdk.NewDecWithPrec(1049, 2) // swap rate
	anotherRandomPrice = sdk.NewDecWithPrec(4882, 2) // swap rate

	oracleDecPrecision = 8
)

func setup(t *testing.T) (testInput, sdk.Handler) {
	input := createTestInput(t)
	h := NewHandler(input.oracleKeeper)

	defaultOracleParams := DefaultParams()
	defaultOracleParams.VotePeriod = int64(1) // Set to one block for convinience
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	return input, h
}

type testInput struct {
	ctx          sdk.Context
	cdc          *codec.Codec
	accKeeper    auth.AccountKeeper
	bankKeeper   bank.Keeper
	oracleKeeper Keeper
	valset       mock.MockValset
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

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
	keyOracle := sdk.NewKVStoreKey(StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewKVStoreKey(staking.TStoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeIAVL, db)

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

	valset := mock.NewMockValSet()
	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{
			sdk.NewCoin(assets.MicroLunaDenom, uLunaAmt),
			sdk.NewCoin(assets.MicroSDRDenom, uSDRAmt),
		})

		require.NoError(t, err)

		// Add validators
		validator := mock.NewMockValidator(sdk.ValAddress(addr.Bytes()), uLunaAmt)
		valset.Validators = append(valset.Validators, validator)
	}

	oracleKeeper := NewKeeper(
		cdc, keyOracle, valset,
		paramsKeeper.Subspace(DefaultParamspace),
	)

	return testInput{ctx, cdc, accKeeper, bankKeeper, oracleKeeper, valset}
}
