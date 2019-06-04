package mint

import (
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/cosmos/cosmos-sdk/x/staking"

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

	uSDRAmount = sdk.NewInt(1005).MulRaw(assets.MicroUnit)
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
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeTransient, db)

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

	stakingKeeper.SetPool(ctx, staking.InitialPool())
	stakingKeeper.SetParams(ctx, staking.DefaultParams())

	mintKeeper := NewKeeper(
		cdc,
		keyMint,
		stakingKeeper,
		bankKeeper,
		accKeeper,
	)

	for _, addr := range addrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(assets.MicroSDRDenom, uSDRAmount)})
		require.NoError(t, err)
	}

	return testInput{ctx, accKeeper, bankKeeper, mintKeeper}
}

func TestKeeperIssuance(t *testing.T) {
	input := createTestInput(t)
	curDay := sdk.ZeroInt()

	// Should be able to claim genesis issunace
	issuance := input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, uSDRAmount.MulRaw(3), issuance)

	// Lowering issuance works
	err := input.mintKeeper.changeIssuance(input.ctx, assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit).Neg())
	require.Nil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, uSDRAmount.MulRaw(3).Sub(sdk.OneInt().MulRaw(assets.MicroUnit)), issuance)

	// ... but not too much
	err = input.mintKeeper.changeIssuance(input.ctx, assets.MicroSDRDenom, sdk.NewInt(5000).MulRaw(assets.MicroUnit).Neg())
	require.NotNil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, uSDRAmount.MulRaw(3).Sub(sdk.OneInt().MulRaw(assets.MicroUnit)), issuance)

	// Raising issuance works, too
	err = input.mintKeeper.changeIssuance(input.ctx, assets.MicroSDRDenom, sdk.NewInt(986).MulRaw(assets.MicroUnit))
	require.Nil(t, err)
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, sdk.NewInt(4000).MulRaw(assets.MicroUnit), issuance)

	// Moving up one epoch inherits the issuance of previous day
	curDay = curDay.Add(sdk.OneInt())
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, sdk.NewInt(4000).MulRaw(assets.MicroUnit), issuance)

	// ... Even when you move many days
	curDay = curDay.Add(sdk.NewInt(10))
	issuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, sdk.NewInt(4000).MulRaw(assets.MicroUnit), issuance)
}

func TestKeeperMintBurn(t *testing.T) {
	input := createTestInput(t)
	curDay := sdk.ZeroInt()
	issuance := input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)

	// Minting new coins results in an issuance increase
	increment := sdk.NewInt(10).MulRaw(assets.MicroUnit)
	err := input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroSDRDenom, increment))
	require.Nil(t, err)
	newIssuance := input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, issuance.Add(increment), newIssuance)

	// Burning new coins results in an issuance decrease
	decrement := sdk.NewInt(10).MulRaw(assets.MicroUnit)
	err = input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.MicroSDRDenom, decrement))
	require.Nil(t, err)
	newIssuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, issuance, newIssuance)

	// Burning new coins errors if requested to burn too much
	decrement = sdk.NewInt(100000).MulRaw(assets.MicroUnit)
	err = input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.MicroSDRDenom, decrement))
	require.NotNil(t, err)
	newIssuance = input.mintKeeper.GetIssuance(input.ctx, assets.MicroSDRDenom, curDay)
	require.Equal(t, issuance, newIssuance)
}

func TestKeeperSeigniorage(t *testing.T) {
	input := createTestInput(t)

	for e := 0; e < 3; e++ {
		input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch * int64(e))
		for i := 0; i < 100; i++ {
			input.mintKeeper.AddSeigniorage(input.ctx, sdk.NewInt(int64(10*(e+1))))
		}
	}

	require.Equal(t, sdk.NewInt(1000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(0)))
	require.Equal(t, sdk.NewInt(2000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(1)))
	require.Equal(t, sdk.NewInt(3000), input.mintKeeper.PeekSeignioragePool(input.ctx, sdk.NewInt(2)))
}
