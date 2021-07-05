package keeper

//nolint
//DONTCOVER

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/ibc-go/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/modules/core"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	customauth "github.com/terra-money/core/custom/auth"
	custombank "github.com/terra-money/core/custom/bank"
	bankwasm "github.com/terra-money/core/custom/bank/wasm"
	customdistr "github.com/terra-money/core/custom/distribution"
	distrwasm "github.com/terra-money/core/custom/distribution/wasm"
	govwasm "github.com/terra-money/core/custom/gov/wasm"
	customparams "github.com/terra-money/core/custom/params"
	customstaking "github.com/terra-money/core/custom/staking"
	stakingwasm "github.com/terra-money/core/custom/staking/wasm"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market"
	marketkeeper "github.com/terra-money/core/x/market/keeper"
	markettypes "github.com/terra-money/core/x/market/types"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	"github.com/terra-money/core/x/oracle"
	oraclekeeper "github.com/terra-money/core/x/oracle/keeper"
	oracletypes "github.com/terra-money/core/x/oracle/types"
	oraclewasm "github.com/terra-money/core/x/oracle/wasm"
	treasurykeeper "github.com/terra-money/core/x/treasury/keeper"
	treasurytypes "github.com/terra-money/core/x/treasury/types"
	treasurywasm "github.com/terra-money/core/x/treasury/wasm"
	"github.com/terra-money/core/x/wasm/config"
	"github.com/terra-money/core/x/wasm/types"
)

const faucetAccountName = "faucet"

// ModuleBasics nolint
var ModuleBasics = module.NewBasicManager(
	customauth.AppModuleBasic{},
	custombank.AppModuleBasic{},
	customstaking.AppModuleBasic{},
	customdistr.AppModuleBasic{},
	customparams.AppModuleBasic{},
	oracle.AppModuleBasic{},
	market.AppModuleBasic{},
	ibc.AppModuleBasic{},
	transfer.AppModuleBasic{},
	capability.AppModuleBasic{},
)

// MakeTestCodec nolint
func MakeTestCodec(t *testing.T) codec.Codec {
	return MakeEncodingConfig(t).Marshaler
}

// MakeEncodingConfig nolint
func MakeEncodingConfig(_ *testing.T) simparams.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)
	types.RegisterInterfaces(interfaceRegistry)

	return simparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// Test Account
var (
	valPubKeys = simapp.CreateTestPubKeys(5)

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

// TestInput nolint
type TestInput struct {
	Ctx                sdk.Context
	Cdc                *codec.LegacyAmino
	AccKeeper          authkeeper.AccountKeeper
	BankKeeper         bankkeeper.Keeper
	StakingKeeper      stakingkeeper.Keeper
	DistributionKeeper distrkeeper.Keeper
	OracleKeeper       oraclekeeper.Keeper
	MarketKeeper       marketkeeper.Keeper
	TreasuryKeeper     treasurykeeper.Keeper
	WasmKeeper         Keeper
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	tempDir := t.TempDir()

	keyContract := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distrtypes.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracletypes.StoreKey)
	keyMarket := sdk.NewKVStoreKey(markettypes.StoreKey)
	keyTreasury := sdk.NewKVStoreKey(treasurytypes.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())
	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	ms.MountStoreWithDB(keyContract, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMarket, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTreasury, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		faucetAccountName:              true,
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
		markettypes.ModuleName:         true,
		treasurytypes.ModuleName:       true,
	}

	maccPerms := map[string][]string{
		faucetAccountName:              {authtypes.Burner, authtypes.Minter},
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          nil,
		oracletypes.ModuleName:         nil,
		markettypes.ModuleName:         {authtypes.Burner, authtypes.Minter},
		treasurytypes.ModuleName:       {authtypes.Minter},
		treasurytypes.BurnModuleName:   {authtypes.Burner},
	}

	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, keyParams, tkeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(appCodec, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), blackListAddrs)

	totalSupply := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)*10))))
	err := bankKeeper.MintCoins(ctx, faucetAccountName, totalSupply)
	require.NoError(t, err)

	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

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
		authtypes.FeeCollectorName, blackListAddrs)

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
	marketAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Burner, authtypes.Minter)

	err = bankKeeper.SendCoinsFromModuleToModule(ctx, faucetAccountName, stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs))))))
	require.NoError(t, err)

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)
	accountKeeper.SetModuleAccount(ctx, oracleAcc)
	accountKeeper.SetModuleAccount(ctx, marketAcc)

	for _, addr := range Addrs {
		accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(addr))
		err := bankKeeper.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr, InitCoins)
		require.NoError(t, err)
	}

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
		accountKeeper, bankKeeper, oracleKeeper,
	)
	marketKeeper.SetParams(ctx, markettypes.DefaultParams())

	treasuryKeeper := treasurykeeper.NewKeeper(
		appCodec,
		keyTreasury, paramsKeeper.Subspace(treasurytypes.ModuleName),
		accountKeeper, bankKeeper,
		marketKeeper, oracleKeeper,
		stakingKeeper, distrKeeper,
		distrtypes.ModuleName,
	)

	treasuryKeeper.SetParams(ctx, treasurytypes.DefaultParams())

	router := baseapp.NewMsgServiceRouter()
	querier := baseapp.NewGRPCQueryRouter()
	banktypes.RegisterQueryServer(querier, bankKeeper)
	stakingtypes.RegisterQueryServer(querier, stakingkeeper.Querier{Keeper: stakingKeeper})
	distrtypes.RegisterQueryServer(querier, distrKeeper)

	keeper := NewKeeper(
		appCodec,
		keyContract,
		paramsKeeper.Subspace(types.ModuleName),
		accountKeeper,
		bankKeeper,
		treasuryKeeper,
		router,
		querier,
		types.DefaultFeatures,
		tempDir,
		config.DefaultConfig(),
	)

	router.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	bankMsgServer := bankkeeper.NewMsgServerImpl(bankKeeper)
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)
	distrMsgServer := distrkeeper.NewMsgServerImpl(distrKeeper)
	marketMsgServer := marketkeeper.NewMsgServerImpl(marketKeeper)
	wasmMsgServer := NewMsgServerImpl(keeper)

	banktypes.RegisterMsgServer(router, bankMsgServer)
	stakingtypes.RegisterMsgServer(router, stakingMsgServer)
	distrtypes.RegisterMsgServer(router, distrMsgServer)
	markettypes.RegisterMsgServer(router, marketMsgServer)
	types.RegisterMsgServer(router, wasmMsgServer)

	keeper.SetParams(ctx, types.DefaultParams())
	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		types.WasmQueryRouteBank:     bankwasm.NewWasmQuerier(bankKeeper),
		types.WasmQueryRouteStaking:  stakingwasm.NewWasmQuerier(stakingKeeper, distrKeeper),
		types.WasmQueryRouteMarket:   marketwasm.NewWasmQuerier(marketKeeper),
		types.WasmQueryRouteTreasury: treasurywasm.NewWasmQuerier(treasuryKeeper),
		types.WasmQueryRouteWasm:     NewWasmQuerier(keeper),
		types.WasmQueryRouteOracle:   oraclewasm.NewWasmQuerier(oracleKeeper),
	}, NewStargateWasmQuerier(keeper))
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		types.WasmMsgParserRouteBank:         bankwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteStaking:      stakingwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteMarket:       marketwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteDistribution: distrwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteGov:          govwasm.NewWasmMsgParser(),
		types.WasmMsgParserRouteWasm:         NewWasmMsgParser(),
	}, NewStargateWasmMsgParser(legacyAmino))

	keeper.SetLastCodeID(ctx, 0)
	keeper.SetLastInstanceID(ctx, 0)

	return TestInput{
		ctx.WithGasMeter(sdk.NewGasMeter(keeper.MaxContractGas(ctx))),
		legacyAmino,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distrKeeper,
		oracleKeeper,
		marketKeeper,
		treasuryKeeper,
		keeper}
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

func createFakeFundedAccount(ctx sdk.Context, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	ak.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(addr))

	if err := bk.MintCoins(ctx, faucetAccountName, coins); err != nil {
		panic(err)
	}

	if err := bk.SendCoinsFromModuleToAccount(ctx, faucetAccountName, addr, coins); err != nil {
		panic(err)
	}
	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// HackatomExampleInitMsg nolint
type HackatomExampleInitMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}
