package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	customauth "github.com/classic-terra/core/v2/custom/auth"
	custombank "github.com/classic-terra/core/v2/custom/bank"
	customdistr "github.com/classic-terra/core/v2/custom/distribution"
	customparams "github.com/classic-terra/core/v2/custom/params"
	customstaking "github.com/classic-terra/core/v2/custom/staking"
	core "github.com/classic-terra/core/v2/types"
	"github.com/classic-terra/core/v2/x/market"
	marketkeeper "github.com/classic-terra/core/v2/x/market/keeper"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
	"github.com/classic-terra/core/v2/x/oracle"
	oraclekeeper "github.com/classic-terra/core/v2/x/oracle/keeper"
	oracletypes "github.com/classic-terra/core/v2/x/oracle/types"
	"github.com/classic-terra/core/v2/x/treasury/types"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const faucetAccountName = "faucet"

var ModuleBasics = module.NewBasicManager(
	customauth.AppModuleBasic{},
	custombank.AppModuleBasic{},
	customdistr.AppModuleBasic{},
	customstaking.AppModuleBasic{},
	customparams.AppModuleBasic{},
	oracle.AppModuleBasic{},
	market.AppModuleBasic{},
)

func MakeTestCodec(t *testing.T) codec.Codec {
	return MakeEncodingConfig(t).Codec
}

func MakeEncodingConfig(_ *testing.T) simparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return simparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

var (
	ValPubKeys = simapp.CreateTestPubKeys(5)

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

	InitTokens = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens))
)

type TestInput struct {
	Ctx            sdk.Context
	Cdc            *codec.LegacyAmino
	TreasuryKeeper Keeper
	AccountKeeper  authkeeper.AccountKeeper
	BankKeeper     bankkeeper.Keeper
	DistrKeeper    distrkeeper.Keeper
	StakingKeeper  stakingkeeper.Keeper
	MarketKeeper   types.MarketKeeper
	OracleKeeper   types.OracleKeeper
}

func CreateTestInput(t *testing.T) TestInput {
	sdk.GetConfig().SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)

	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(oracletypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distrtypes.StoreKey)
	keyMarket := sdk.NewKVStoreKey(markettypes.StoreKey)
	keyTreasury := sdk.NewKVStoreKey(types.StoreKey)
	keyWasm := sdk.NewKVStoreKey(wasmtypes.StoreKey)
	// keyIbcHost := sdk.NewKVStoreKey(ibchost.StoreKey)
	keyCapability := sdk.NewKVStoreKey(capabilitytypes.StoreKey)
	// keyUpgrade := sdk.NewKVStoreKey(upgradetypes.StoreKey)
	memKeyCapability := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Codec, encodingConfig.Amino

	ms.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, storetypes.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, storetypes.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
		oracletypes.ModuleName:         true,
		faucetAccountName:              true,
	}

	maccPerms := map[string][]string{
		faucetAccountName:              {authtypes.Minter, authtypes.Burner},
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		markettypes.ModuleName:         {authtypes.Burner, authtypes.Minter},
		distrtypes.ModuleName:          nil,
		oracletypes.ModuleName:         nil,
		types.ModuleName:               {authtypes.Burner, authtypes.Minter},
		types.BurnModuleName:           {authtypes.Burner},
	}

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tKeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms, sdk.GetConfig().GetBech32AccountAddrPrefix())
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), blackListAddrs)

	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)*10))))
	bankKeeper.MintCoins(ctx, faucetAccountName, totalSupply)

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
		accountKeeper, bankKeeper, stakingKeeper,
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
	oracleAcc := authtypes.NewEmptyModuleAccount(oracletypes.ModuleName)
	marketAcc := authtypes.NewEmptyModuleAccount(markettypes.ModuleName, authtypes.Burner, authtypes.Minter)
	treasuryAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Burner, authtypes.Minter)
	burnAcc := authtypes.NewEmptyModuleAccount(types.BurnModuleName, authtypes.Burner)

	// + 1 for burn account
	bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)+1)))))

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)
	accountKeeper.SetModuleAccount(ctx, oracleAcc)
	accountKeeper.SetModuleAccount(ctx, marketAcc)
	accountKeeper.SetModuleAccount(ctx, treasuryAcc)
	accountKeeper.SetModuleAccount(ctx, burnAcc)

	for _, addr := range Addrs {
		accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(addr))
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr, InitCoins)
		require.NoError(t, err)
	}

	// to test burn module account
	err := bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, types.BurnModuleName, InitCoins)
	require.NoError(t, err)

	capabilityKeeper := capabilitykeeper.NewKeeper(
		appCodec, keyCapability, memKeyCapability[capabilitytypes.MemStoreKey],
	)

	// mock wasm
	scopedWasmKeeper := capabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	wasmConfig := wasmtypes.DefaultWasmConfig()
	supportedFeatures := "iterator,staking,stargate,terra,cosmwasm_1_1"
	wasmOpts := []wasmkeeper.Option{}
	wasmKeeper := wasmkeeper.NewKeeper(
		appCodec, keyWasm,
		paramsKeeper.Subspace(wasmtypes.ModuleName),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distrKeeper,
		nil,
		nil,
		scopedWasmKeeper,
		nil,
		nil,
		nil,
		"",
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	oracleKeeper := oraclekeeper.NewKeeper(
		appCodec,
		keyOracle,
		paramsKeeper.Subspace(oracletypes.ModuleName),
		accountKeeper,
		bankKeeper,
		distrKeeper,
		stakingKeeper,
		distrtypes.ModuleName,
	)
	oracleDefaultParams := oracletypes.DefaultParams()
	oracleKeeper.SetParams(ctx, oracleDefaultParams)

	for _, denom := range oracleDefaultParams.Whitelist {
		oracleKeeper.SetTobinTax(ctx, denom.Name, denom.TobinTax)
	}

	marketKeeper := marketkeeper.NewKeeper(
		appCodec,
		keyMarket, paramsKeeper.Subspace(markettypes.ModuleName),
		accountKeeper,
		bankKeeper,
		oracleKeeper,
	)
	marketKeeper.SetParams(ctx, markettypes.DefaultParams())

	treasuryKeeper := NewKeeper(
		appCodec,
		keyTreasury, paramsKeeper.Subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		marketKeeper,
		oracleKeeper,
		stakingKeeper,
		distrKeeper,
		&wasmKeeper,
		distrtypes.ModuleName,
	)

	treasuryKeeper.SetParams(ctx, types.DefaultParams())

	return TestInput{ctx, legacyAmino, treasuryKeeper, accountKeeper, bankKeeper, distrKeeper, stakingKeeper, marketKeeper, oracleKeeper}
}

// NewTestMsgCreateValidator test msg creator
func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey cryptotypes.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	msg, _ := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(core.MicroLunaDenom, amt),
		stakingtypes.Description{}, commission, sdk.OneInt(),
	)

	return msg
}

func setupValidators(t *testing.T) (TestInput, stakingtypes.MsgServer) {
	input := CreateTestInput(t)
	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(input.StakingKeeper)

	// Create Validators
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	_, err := stakingMsgSvr.CreateValidator(input.Ctx, NewTestMsgCreateValidator(addr, val, amt))

	require.NoError(t, err)
	_, err = stakingMsgSvr.CreateValidator(input.Ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)

	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, stakingMsgSvr
}

// FundAccount is a utility function that funds an account by minting and
// sending the coins to the address. This should be used for testing purposes
// only!
func FundAccount(input TestInput, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := input.BankKeeper.MintCoins(input.Ctx, faucetAccountName, amounts); err != nil {
		return err
	}

	return input.BankKeeper.SendCoinsFromModuleToAccount(input.Ctx, faucetAccountName, addr, amounts)
}
