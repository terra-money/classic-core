package keepers

import (
	"path/filepath"

	icahostkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	customwasmkeeper "github.com/classic-terra/core/v2/custom/wasm/keeper"
	terrawasm "github.com/classic-terra/core/v2/wasmbinding"

	feesharekeeper "github.com/classic-terra/core/v2/x/feeshare/keeper"
	feesharetypes "github.com/classic-terra/core/v2/x/feeshare/types"
	marketkeeper "github.com/classic-terra/core/v2/x/market/keeper"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
	oraclekeeper "github.com/classic-terra/core/v2/x/oracle/keeper"
	oracletypes "github.com/classic-terra/core/v2/x/oracle/types"
	treasurykeeper "github.com/classic-terra/core/v2/x/treasury/keeper"
	treasurytypes "github.com/classic-terra/core/v2/x/treasury/types"
)

type AppKeepers struct {
	// appKeepers.keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	AuthzKeeper      authzkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the appKeepers, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	OracleKeeper     oraclekeeper.Keeper
	FeeShareKeeper   feesharekeeper.Keeper
	MarketKeeper     marketkeeper.Keeper
	TreasuryKeeper   treasurykeeper.Keeper
	WasmKeeper       wasmkeeper.Keeper

	ICAHostKeeper icahostkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
}

func NewAppKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	allowedReceivingModAcc map[string]bool,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	wasmOpts []wasm.Option,
	appOpts servertypes.AppOptions,
) AppKeepers {
	appKeepers := AppKeepers{}

	appKeepers.GenerateKeys()

	// init params keeper and subspaces
	appKeepers.ParamsKeeper = initParamsKeeper(appCodec, cdc, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	appKeepers.CapabilityKeeper.Seal()

	// add keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, appKeepers.keys[authtypes.StoreKey], appKeepers.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, appKeepers.keys[banktypes.StoreKey], appKeepers.AccountKeeper, appKeepers.GetSubspace(banktypes.ModuleName), appKeepers.BlacklistedAccAddrs(maccPerms, allowedReceivingModAcc),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, appKeepers.keys[stakingtypes.StoreKey], appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.GetSubspace(stakingtypes.ModuleName),
	)
	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec, appKeepers.keys[minttypes.StoreKey], appKeepers.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper, authtypes.FeeCollectorName,
	)
	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, appKeepers.keys[distrtypes.StoreKey], appKeepers.GetSubspace(distrtypes.ModuleName), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, appKeepers.BlacklistedAccAddrs(maccPerms, allowedReceivingModAcc),
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, appKeepers.keys[slashingtypes.StoreKey], &stakingKeeper, appKeepers.GetSubspace(slashingtypes.ModuleName),
	)
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName), invCheckPeriod, appKeepers.BankKeeper, authtypes.FeeCollectorName,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, appKeepers.keys[feegrant.StoreKey], appKeepers.AccountKeeper)
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, appKeepers.keys[upgradetypes.StoreKey], appCodec, homePath, bApp)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistrKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(appKeepers.keys[authzkeeper.StoreKey], appCodec, bApp.MsgServiceRouter())

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, appKeepers.keys[ibchost.StoreKey], appKeepers.GetSubspace(ibchost.ModuleName), appKeepers.StakingKeeper, appKeepers.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, appKeepers.keys[ibctransfertypes.StoreKey], appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.ChannelKeeper, &appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper, scopedTransferKeeper,
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := appKeepers.getIBCRouter()
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, appKeepers.keys[evidencetypes.StoreKey], &appKeepers.StakingKeeper, appKeepers.SlashingKeeper,
	)
	// If evidence needs to be handled for the appKeepers, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	// Initialize terra module keepers
	appKeepers.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec, appKeepers.keys[oracletypes.StoreKey], appKeepers.GetSubspace(oracletypes.ModuleName),
		appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.DistrKeeper, &stakingKeeper, distrtypes.ModuleName,
	)
	appKeepers.MarketKeeper = marketkeeper.NewKeeper(
		appCodec, appKeepers.keys[markettypes.StoreKey],
		appKeepers.GetSubspace(markettypes.ModuleName),
		appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.OracleKeeper,
	)
	appKeepers.TreasuryKeeper = treasurykeeper.NewKeeper(
		appCodec, appKeepers.keys[treasurytypes.StoreKey],
		appKeepers.GetSubspace(treasurytypes.ModuleName),
		appKeepers.AccountKeeper, appKeepers.BankKeeper,
		appKeepers.MarketKeeper, appKeepers.OracleKeeper,
		appKeepers.StakingKeeper, appKeepers.DistrKeeper,
		distrtypes.ModuleName)

	appKeepers.FeeShareKeeper = feesharekeeper.NewKeeper(
		appKeepers.keys[feesharetypes.StoreKey],
		appCodec,
		appKeepers.GetSubspace(feesharetypes.ModuleName),
		appKeepers.BankKeeper,
		appKeepers.WasmKeeper,
		authtypes.FeeCollectorName,
	)

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	supportedFeatures := "iterator,staking,stargate,terra,cosmwasm_1_1"

	wasmMsgHandler := customwasmkeeper.NewMessageHandler(bApp.MsgServiceRouter(), appKeepers.IBCKeeper.ChannelKeeper, scopedWasmKeeper, appKeepers.BankKeeper, appKeepers.TreasuryKeeper, appKeepers.AccountKeeper, appCodec, appKeepers.TransferKeeper)
	// the first slice will replace all default msh handler with custom one
	wasmOpts = append([]wasmkeeper.Option{wasmkeeper.WithMessageHandler(wasmMsgHandler)}, wasmOpts...)
	// the second slice will add custom querier and message handler decorator
	// this order must be uphold else error will be thrown
	wasmOpts = append(wasmOpts, terrawasm.RegisterCustomPlugins(&appKeepers.MarketKeeper, &appKeepers.OracleKeeper, &appKeepers.TreasuryKeeper)...)

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[wasm.StoreKey],
		appKeepers.GetSubspace(wasm.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		filepath.Join(homePath, "data"),
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	// register the proposal types
	govRouter := appKeepers.getGovRouter()
	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec, appKeepers.keys[govtypes.StoreKey], appKeepers.GetSubspace(govtypes.ModuleName), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		&stakingKeeper, govRouter,
	)

	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedICAHostKeeper = scopedICAHostKeeper
	appKeepers.ScopedTransferKeeper = scopedTransferKeeper

	return appKeepers
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(markettypes.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)
	paramsKeeper.Subspace(treasurytypes.ModuleName)
	paramsKeeper.Subspace(feesharetypes.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (appKeepers *AppKeepers) BlacklistedAccAddrs(maccPerms map[string][]string, allowedReceivingModAcc map[string]bool) map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}
