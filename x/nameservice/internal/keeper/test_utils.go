package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/oracle"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/terra-project/core/x/nameservice/internal/types"
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

	InitTokens = sdk.TokensFromConsensusPower(20000)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens), sdk.NewCoin(core.MicroSDRDenom, InitTokens))
)

// TestInput nolint
type TestInput struct {
	Ctx               sdk.Context
	Cdc               *codec.Codec
	AccKeeper         auth.AccountKeeper
	BankKeeper        bank.Keeper
	SupplyKeeper      supply.Keeper
	MarketKeeper      market.Keeper
	OracleKeeper      oracle.Keeper
	NameserviceKeeper Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
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
	tKeyStaking := sdk.NewKVStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyMarket := sdk.NewKVStoreKey(market.StoreKey)
	keyNameservice := sdk.NewKVStoreKey(types.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyNameservice, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName:     true,
		staking.NotBondedPoolName: true,
		staking.BondedPoolName:    true,
		distr.ModuleName:          true,
		oracle.ModuleName:         true,
		market.ModuleName:         true,
		types.ModuleName:          true,
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blackListAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		distr.ModuleName:          nil,
		market.ModuleName:         {supply.Burner, supply.Minter},
		oracle.ModuleName:         nil,
		types.ModuleName:          {supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)))),
		sdk.NewCoin(core.MicroSDRDenom, InitTokens.MulRaw(int64(len(Addrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking, tKeyStaking,
		supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	distrKeeper := distr.NewKeeper(
		cdc,
		keyDistr, paramsKeeper.Subspace(distr.DefaultParamspace),
		stakingKeeper, supplyKeeper, distr.DefaultCodespace,
		auth.FeeCollectorName, blackListAddrs,
	)

	distrKeeper.SetFeePool(ctx, distr.InitialFeePool())
	distrKeeper.SetCommunityTax(ctx, sdk.NewDecWithPrec(2, 2))
	distrKeeper.SetBaseProposerReward(ctx, sdk.NewDecWithPrec(1, 2))
	distrKeeper.SetBonusProposerReward(ctx, sdk.NewDecWithPrec(4, 2))

	oracleKeeper := oracle.NewKeeper(
		cdc,
		keyOracle, paramsKeeper.Subspace(oracle.DefaultParamspace),
		distrKeeper, stakingKeeper, supplyKeeper, distr.ModuleName,
		oracle.DefaultCodespace,
	)

	marketKeeper := market.NewKeeper(
		cdc,
		keyMarket, paramsKeeper.Subspace(market.DefaultParamspace),
		oracleKeeper, supplyKeeper, market.DefaultCodespace,
	)

	nameserviceKeeper := NewKeeper(
		cdc,
		keyNameservice, paramsKeeper.Subspace(types.DefaultParamspace),
		supplyKeeper, marketKeeper, types.DefaultCodespace,
	)

	nameserviceKeeper.SetParams(ctx, types.DefaultParams())

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	distrAcc := supply.NewEmptyModuleAccount(distr.ModuleName)
	marketAcc := supply.NewEmptyModuleAccount(market.ModuleName, supply.Burner, supply.Minter)
	oracleAcc := supply.NewEmptyModuleAccount(oracle.ModuleName, supply.Minter)

	_ = notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs))))))

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, distrAcc)
	supplyKeeper.SetModuleAccount(ctx, marketAcc)
	supplyKeeper.SetModuleAccount(ctx, oracleAcc)

	genesis := staking.DefaultGenesisState()
	genesis.Params.BondDenom = core.MicroLunaDenom
	_ = staking.InitGenesis(ctx, stakingKeeper, accountKeeper, supplyKeeper, genesis)

	for _, addr := range Addrs {
		_, err := bankKeeper.AddCoins(ctx, addr, InitCoins)
		require.NoError(t, err)
	}

	stakingKeeper.SetHooks(staking.NewMultiStakingHooks(distrKeeper.Hooks()))

	return TestInput{Ctx: ctx, Cdc: cdc, AccKeeper: accountKeeper, BankKeeper: bankKeeper, NameserviceKeeper: nameserviceKeeper, SupplyKeeper: supplyKeeper, MarketKeeper: marketKeeper, OracleKeeper: oracleKeeper}
}
