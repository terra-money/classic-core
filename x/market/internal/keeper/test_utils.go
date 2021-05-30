// nolint:deadcode unused noalias
package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/internal/types"
	"github.com/terra-money/core/x/oracle"

	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var (
	PubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(PubKeys[0].Address()),
		sdk.AccAddress(PubKeys[1].Address()),
		sdk.AccAddress(PubKeys[2].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(PubKeys[0].Address()),
		sdk.ValAddress(PubKeys[1].Address()),
		sdk.ValAddress(PubKeys[2].Address()),
	}

	InitTokens = sdk.TokensFromConsensusPower(200)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens))
)

// TestInput nolint
type TestInput struct {
	Ctx          sdk.Context
	Cdc          *codec.Codec
	Acckeeper    auth.AccountKeeper
	OracleKeeper oracle.Keeper
	SupplyKeeper supply.Keeper
	MarketKeeper Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	oracle.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	params.RegisterCodec(cdc)

	return cdc
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyMarket := sdk.NewKVStoreKey(types.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName:     true,
		staking.NotBondedPoolName: true,
		staking.BondedPoolName:    true,
		distr.ModuleName:          true,
		types.ModuleName:          true,
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), blackListAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		distr.ModuleName:          nil,
		oracle.ModuleName:         nil,
		types.ModuleName:          {supply.Burner, supply.Minter},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(
		cdc, keyStaking, supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
	)

	distrKeeper := distr.NewKeeper(
		cdc, keyDistr, paramsKeeper.Subspace(distr.DefaultParamspace),
		stakingKeeper, supplyKeeper, auth.FeeCollectorName, blackListAddrs,
	)

	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle, paramsKeeper.Subspace(oracle.DefaultParamspace),
		distrKeeper, stakingKeeper, supplyKeeper, distr.ModuleName,
	)

	keeper := NewKeeper(
		cdc,
		keyMarket, paramsKeeper.Subspace(types.DefaultParamspace),
		oracleKeeper, supplyKeeper,
	)

	keeper.SetParams(ctx, types.DefaultParams())
	oracleKeeper.SetParams(ctx, oracle.DefaultParams())

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	distrAcc := supply.NewEmptyModuleAccount(distr.ModuleName)
	marketAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner, supply.Minter)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, distrAcc)
	supplyKeeper.SetModuleAccount(ctx, marketAcc)

	for _, addr := range Addrs {
		_, err := bankKeeper.AddCoins(ctx, sdk.AccAddress(addr), InitCoins)
		require.NoError(t, err)
	}

	supply := supplyKeeper.GetSupply(ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs))))))
	supplyKeeper.SetSupply(ctx, supply)

	return TestInput{ctx, cdc, accountKeeper, oracleKeeper, supplyKeeper, keeper}
}
