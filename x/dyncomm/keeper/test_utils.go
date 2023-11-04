package keeper

//nolint
//DONTCOVER

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	customauth "github.com/classic-terra/core/v2/custom/auth"
	custombank "github.com/classic-terra/core/v2/custom/bank"
	customdistr "github.com/classic-terra/core/v2/custom/distribution"
	customparams "github.com/classic-terra/core/v2/custom/params"
	customstaking "github.com/classic-terra/core/v2/custom/staking"
	core "github.com/classic-terra/core/v2/types"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	types "github.com/classic-terra/core/v2/x/dyncomm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const faucetAccountName = "faucet"

var ModuleBasics = module.NewBasicManager(
	customauth.AppModuleBasic{},
	custombank.AppModuleBasic{},
	customstaking.AppModuleBasic{},
	customdistr.AppModuleBasic{},
	customparams.AppModuleBasic{},
)

// MakeTestCodec
func MakeTestCodec(t *testing.T) codec.Codec {
	return MakeEncodingConfig(t).Codec
}

// MakeEncodingConfig
func MakeEncodingConfig(_ *testing.T) simparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)
	types.RegisterInterfaces(interfaceRegistry)

	return simparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// Test Account
var (
	PubKeys = simapp.CreateTestPubKeys(32)

	InitTokens    = sdk.TokensFromConsensusPower(10_000, sdk.DefaultPowerReduction)
	InitCoins     = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens))
	DelegateCoins = sdk.NewCoin(core.MicroLunaDenom, InitTokens)

	blackListAddrs = map[string]bool{
		faucetAccountName:              true,
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
	}

	maccPerms = map[string][]string{
		faucetAccountName:              {authtypes.Minter},
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          nil,
		types.ModuleName:               {authtypes.Burner, authtypes.Minter},
	}
)

type TestInput struct {
	Ctx           sdk.Context
	Cdc           *codec.LegacyAmino
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	DistrKeeper   distrkeeper.Keeper
	StakingKeeper stakingkeeper.Keeper
	DyncommKeeper Keeper
}

func CreateTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distrtypes.StoreKey)
	keyDyncomm := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Codec, encodingConfig.Amino

	ms.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, storetypes.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDyncomm, storetypes.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tKeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms, sdk.GetConfig().GetBech32AccountAddrPrefix())
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), blackListAddrs)

	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, math.Int(math.LegacyNewDec(1_000_000_000_000))))
	err := bankKeeper.MintCoins(ctx, faucetAccountName, totalSupply)
	require.NoError(t, err)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keyStaking,
		accountKeeper,
		bankKeeper,
		paramsKeeper.Subspace(stakingtypes.ModuleName),
	)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = core.MicroLunaDenom
	stakingKeeper.SetParams(ctx, stakingParams)

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		keyDistr, paramsKeeper.Subspace(distrtypes.ModuleName),
		accountKeeper, bankKeeper, &stakingKeeper,
		authtypes.FeeCollectorName)

	distrKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	distrParams := distrtypes.DefaultParams()
	distrParams.CommunityTax = sdk.NewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	distrKeeper.SetParams(ctx, distrParams)
	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distrKeeper.Hooks()))

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName)

	err = bankKeeper.SendCoinsFromModuleToModule(
		ctx, faucetAccountName, stakingtypes.NotBondedPoolName,
		sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(PubKeys))))),
	)
	require.NoError(t, err)

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)

	for idx := range PubKeys {
		accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(AddrFrom(idx)))
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, AddrFrom(idx), InitCoins)
		require.NoError(t, err)
	}

	dyncommKeeper := NewKeeper(
		appCodec, keyDyncomm,
		paramsKeeper.Subspace(types.ModuleName),
		stakingKeeper,
	)
	dyncommKeeper.SetParams(
		ctx, types.DefaultParams(),
	)

	return TestInput{ctx, legacyAmino, accountKeeper, bankKeeper, distrKeeper, stakingKeeper, dyncommKeeper}
}

func CallCreateValidatorHooks(ctx sdk.Context, k distrkeeper.Keeper, addr sdk.AccAddress, valAddr sdk.ValAddress) error {
	err := k.Hooks().AfterValidatorCreated(ctx, valAddr)
	if err != nil {
		return err
	}
	err = k.Hooks().BeforeDelegationCreated(ctx, addr, valAddr)
	if err != nil {
		return err
	}
	err = k.Hooks().AfterDelegationModified(ctx, addr, valAddr)
	if err != nil {
		return err
	}
	return nil
}

func CreateValidator(idx int, stake math.Int) (stakingtypes.Validator, error) {
	val, err := stakingtypes.NewValidator(
		ValAddrFrom(idx), PubKeys[idx], stakingtypes.Description{Moniker: "TestValidator"},
	)
	val.Tokens = stake
	val.DelegatorShares = sdk.NewDec(val.Tokens.Int64())
	return val, err
}

func GetPubKey(idx int) (sdkcrypto.PubKey, sdk.AccAddress, sdk.ValAddress) {
	addr := AddrFrom(idx)
	valAddr := ValAddrFrom(idx)
	return PubKeys[idx], addr, valAddr
}

func AddrFrom(idx int) sdk.AccAddress {
	return sdk.AccAddress(PubKeys[idx].Address())
}

func ValAddrFrom(idx int) sdk.ValAddress {
	return sdk.ValAddress(PubKeys[idx].Address())
}
