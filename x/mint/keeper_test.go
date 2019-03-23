package mint

import (
	"terra/types/assets"
	"terra/types/util"
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
	ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.Keeper
	mintKeeper Keeper
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
	keyMint := sdk.NewKVStoreKey(StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
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

	mintKeeper := NewKeeper(
		cdc,
		keyMint,
		bankKeeper,
		accKeeper,
	)

	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(assets.SDRDenom, initAmt)})
		require.NoError(t, err)
	}

	return testInput{ctx, accKeeper, bankKeeper, mintKeeper}
}

func TestKeeperIssuance(t *testing.T) {
	input := createTestInput(t)
	curEpoch := util.GetEpoch(input.ctx)

	// Should be able to claim genesis issunace
	issuance := input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, initAmt.MulRaw(3), issuance)

	// Lowering issuance works
	err := input.mintKeeper.changeIssuance(input.ctx, assets.SDRDenom, sdk.OneInt().Neg())
	require.Nil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, initAmt.MulRaw(3).Sub(sdk.OneInt()), issuance)

	// ... but not too much
	err = input.mintKeeper.changeIssuance(input.ctx, assets.SDRDenom, sdk.NewInt(5000).Neg())
	require.NotNil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, initAmt.MulRaw(3).Sub(sdk.OneInt()), issuance)

	// Raising issuance works, too
	err = input.mintKeeper.changeIssuance(input.ctx, assets.SDRDenom, sdk.NewInt(986))
	require.Nil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, sdk.NewInt(4000), issuance)

	// Moving up one epoch inherits the issuance of previous epoch
	curEpoch = curEpoch.Add(sdk.OneInt())
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, sdk.NewInt(4000), issuance)

	// ... Even when you move many epochs
	curEpoch = curEpoch.Add(sdk.NewInt(10))
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, sdk.NewInt(4000), issuance)
}

func TestKeeperMintBurn(t *testing.T) {
	input := createTestInput(t)
	curEpoch := util.GetEpoch(input.ctx)
	issuance := input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)

	// Minting new coins results in an issuance increase
	increment := sdk.NewInt(10)
	err := input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.SDRDenom, increment))
	require.Nil(t, err)
	newIssuance := input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, issuance.Add(increment), newIssuance)

	// Burning new coins results in an issuance decrease
	decrement := sdk.NewInt(10)
	err = input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.SDRDenom, decrement))
	require.Nil(t, err)
	newIssuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, issuance, newIssuance)

	// Burning new coins errors if requested to burn too much
	decrement = sdk.NewInt(100000)
	err = input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.SDRDenom, decrement))
	require.NotNil(t, err)
	newIssuance = input.mintKeeper.GetIssuance(input.ctx, assets.SDRDenom, curEpoch)
	require.Equal(t, issuance, newIssuance)
}

func TestKeeperSeigniorage(t *testing.T) {
	input := createTestInput(t)

	for e := 0; e < 3; e++ {
		input.ctx = input.ctx.WithBlockHeight(util.GetBlocksPerEpoch() * int64(e))
		for i := 0; i < 100; i++ {
			input.mintKeeper.AddSeigniorage(input.ctx, sdk.NewInt(int64(10*(e+1))))
		}
	}

	require.Equal(t, sdk.NewInt(1000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(0)))
	require.Equal(t, sdk.NewInt(2000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(1)))
	require.Equal(t, sdk.NewInt(3000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(2)))
}
