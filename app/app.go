package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	treasuryclient "github.com/terra-money/core/x/treasury/client"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/auth"
	"github.com/terra-money/core/x/auth/ante"
	"github.com/terra-money/core/x/auth/vesting"
	"github.com/terra-money/core/x/bank"
	"github.com/terra-money/core/x/crisis"
	distr "github.com/terra-money/core/x/distribution"
	"github.com/terra-money/core/x/evidence"
	"github.com/terra-money/core/x/genutil"
	"github.com/terra-money/core/x/gov"
	"github.com/terra-money/core/x/market"
	"github.com/terra-money/core/x/mint"
	"github.com/terra-money/core/x/msgauth"
	"github.com/terra-money/core/x/oracle"
	"github.com/terra-money/core/x/params"
	"github.com/terra-money/core/x/slashing"
	"github.com/terra-money/core/x/staking"
	"github.com/terra-money/core/x/supply"
	"github.com/terra-money/core/x/treasury"
	"github.com/terra-money/core/x/upgrade"
	"github.com/terra-money/core/x/wasm"
	wasmconfig "github.com/terra-money/core/x/wasm/config"

	bankwasm "github.com/terra-money/core/x/bank/wasm"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	oraclewasm "github.com/terra-money/core/x/oracle/wasm"
	stakingwasm "github.com/terra-money/core/x/staking/wasm"
	treasurywasm "github.com/terra-money/core/x/treasury/wasm"
)

const appName = "TerraApp"

var (
	// DefaultCLIHome defines default home directories for terracli
	DefaultCLIHome = os.ExpandEnv("$HOME/.terracli")

	// DefaultNodeHome defines default home directories for terrad
	DefaultNodeHome = os.ExpandEnv("$HOME/.terrad")

	// ModuleBasics = The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			treasuryclient.TaxRateUpdateProposalHandler,
			treasuryclient.RewardWeightUpdateProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		evidence.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		oracle.AppModuleBasic{},
		market.AppModuleBasic{},
		treasury.AppModuleBasic{},
		wasm.AppModuleBasic{},
		msgauth.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil, // just added to enable align fee
		bank.BurnModuleName:       {supply.Burner},
		mint.ModuleName:           {supply.Minter},
		market.ModuleName:         {supply.Minter, supply.Burner},
		oracle.ModuleName:         nil,
		distr.ModuleName:          nil,
		treasury.ModuleName:       {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		oracle.ModuleName:   true,
		bank.BurnModuleName: true,
	}
)

// MakeCodec builds application codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc
}

// Verify app interface at compile time
var _ simapp.App = (*TerraApp)(nil)

// TerraApp is Extended ABCI application
type TerraApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	upgradeKeeper  upgrade.Keeper
	evidenceKeeper evidence.Keeper
	oracleKeeper   oracle.Keeper
	marketKeeper   market.Keeper
	mintKeeper     mint.Keeper
	treasuryKeeper treasury.Keeper
	wasmKeeper     wasm.Keeper
	msgauthKeeper  msgauth.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewTerraApp returns a reference to an initialized TerraApp.
func NewTerraApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, wasmConfig *wasmconfig.Config,
	baseAppOptions ...func(*bam.BaseApp)) *TerraApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, oracle.StoreKey,
		market.StoreKey, mint.StoreKey, treasury.StoreKey,
		upgrade.StoreKey, evidence.StoreKey, wasm.StoreKey,
		msgauth.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	var app = &TerraApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[oracle.ModuleName] = app.paramsKeeper.Subspace(oracle.DefaultParamspace)
	app.subspaces[market.ModuleName] = app.paramsKeeper.Subspace(market.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[wasm.ModuleName] = app.paramsKeeper.Subspace(wasm.DefaultParamspace)
	app.subspaces[treasury.ModuleName] = app.paramsKeeper.Subspace(treasury.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(app.cdc, keys[staking.StoreKey], app.supplyKeeper, app.subspaces[staking.ModuleName])
	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(app.cdc, keys[slashing.StoreKey], &stakingKeeper,
		app.subspaces[slashing.ModuleName])
	app.crisisKeeper = crisis.NewKeeper(app.subspaces[crisis.ModuleName], invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)
	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.cdc)
	app.oracleKeeper = oracle.NewKeeper(app.cdc, keys[oracle.StoreKey], app.subspaces[oracle.ModuleName], app.distrKeeper,
		&stakingKeeper, app.supplyKeeper, distr.ModuleName)
	app.marketKeeper = market.NewKeeper(app.cdc, keys[market.StoreKey], app.subspaces[market.ModuleName],
		app.oracleKeeper, app.supplyKeeper)
	app.mintKeeper = mint.NewKeeper(app.cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.treasuryKeeper = treasury.NewKeeper(app.cdc, keys[treasury.StoreKey], app.subspaces[treasury.ModuleName],
		app.supplyKeeper, app.marketKeeper, &stakingKeeper, app.distrKeeper,
		oracle.ModuleName, distr.ModuleName)
	app.msgauthKeeper = msgauth.NewKeeper(app.cdc, keys[msgauth.StoreKey], bApp.Router(),
		bank.MsgSend{}.Type(),
		market.MsgSwap{}.Type(),
		gov.MsgVote{}.Type(),
	)

	// register the evidence router
	evidenceRouter := evidence.NewRouter()
	evidenceKeeper := evidence.NewKeeper(app.cdc, keys[evidence.StoreKey],
		app.subspaces[evidence.ModuleName], &app.stakingKeeper, app.slashingKeeper)
	evidenceKeeper.SetRouter(evidenceRouter)
	app.evidenceKeeper = *evidenceKeeper

	// create wasm keeper with msg parser & querier
	app.wasmKeeper = wasm.NewKeeper(app.cdc, keys[wasm.StoreKey], app.subspaces[wasm.ModuleName],
		app.accountKeeper, app.bankKeeper, app.supplyKeeper, app.treasuryKeeper, bApp.Router(), wasm.DefaultFeatures, wasmConfig)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(treasury.RouterKey, treasury.NewTreasuryPolicyUpdateHandler(app.treasuryKeeper))
	app.govKeeper = gov.NewKeeper(app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName],
		app.supplyKeeper, &stakingKeeper, govRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	app.wasmKeeper.RegisterMsgParsers(map[string]wasm.WasmMsgParserInterface{
		wasm.WasmMsgParserRouteBank:    bankwasm.NewWasmMsgParser(),
		wasm.WasmMsgParserRouteStaking: stakingwasm.NewWasmMsgParser(),
		wasm.WasmMsgParserRouteMarket:  marketwasm.NewWasmMsgParser(),
		wasm.WasmMsgParserRouteWasm:    wasm.NewWasmMsgParser(),
	})
	app.wasmKeeper.RegisterQueriers(map[string]wasm.WasmQuerierInterface{
		wasm.WasmQueryRouteBank:     bankwasm.NewWasmQuerier(app.bankKeeper),
		wasm.WasmQueryRouteStaking:  stakingwasm.NewWasmQuerier(app.stakingKeeper),
		wasm.WasmQueryRouteMarket:   marketwasm.NewWasmQuerier(app.marketKeeper),
		wasm.WasmQueryRouteOracle:   oraclewasm.NewWasmQuerier(app.oracleKeeper),
		wasm.WasmQueryRouteTreasury: treasurywasm.NewWasmQuerier(app.treasuryKeeper),
		wasm.WasmQueryRouteWasm:     wasm.NewWasmQuerier(app.wasmKeeper),
	})

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper, app.supplyKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		market.NewAppModule(app.marketKeeper, app.accountKeeper, app.oracleKeeper),
		oracle.NewAppModule(app.oracleKeeper, app.accountKeeper),
		treasury.NewAppModule(app.treasuryKeeper),
		wasm.NewAppModule(app.wasmKeeper, app.accountKeeper, app.bankKeeper),
		msgauth.NewAppModule(app.msgauthKeeper, app.accountKeeper, app.bankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName,
		evidence.ModuleName, wasm.ModuleName)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, oracle.ModuleName, gov.ModuleName, market.ModuleName,
		treasury.ModuleName, msgauth.ModuleName, staking.ModuleName)

	// genutils must occur after staking so that pools are properly
	// treasury must occur after supply so that initial issuance is properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(auth.ModuleName, distr.ModuleName,
		staking.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, mint.ModuleName, supply.ModuleName,
		oracle.ModuleName, treasury.ModuleName, market.ModuleName,
		wasm.ModuleName, msgauth.ModuleName, crisis.ModuleName,
		genutil.ModuleName, evidence.ModuleName)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// fuzz test simulation
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper, app.supplyKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		params.NewAppModule(), // NOTE: only used for simulation to generate randomized param change proposals
		market.NewAppModule(app.marketKeeper, app.accountKeeper, app.oracleKeeper),
		oracle.NewAppModule(app.oracleKeeper, app.accountKeeper),
		treasury.NewAppModule(app.treasuryKeeper),
		wasm.NewAppModule(app.wasmKeeper, app.accountKeeper, app.bankKeeper),
		msgauth.NewAppModule(app.msgauthKeeper, app.accountKeeper, app.bankKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(
		app.accountKeeper,
		app.supplyKeeper,
		app.oracleKeeper,
		app.treasuryKeeper,
		auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// Name returns the name of the App
func (app *TerraApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker defines application updates at every begin block
func (app *TerraApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker defines application updates at every end block
func (app *TerraApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.mm.EndBlock(ctx, req)
	if core.IsSoftforkHeight(ctx, 3) {
		return abci.ResponseEndBlock{
			ConsensusParamUpdates: &abci.ConsensusParams{
				Block: &abci.BlockParams{
					MaxBytes: 1000000,
					MaxGas:   30000000,
				},
			},
			ValidatorUpdates: res.ValidatorUpdates,
			Events:           res.Events,
		}
	}

	return res
}

// InitChainer defines application update at chain initialization
func (app *TerraApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height
func (app *TerraApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *TerraApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *TerraApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[supply.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns TerraApp's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *TerraApp) Codec() *codec.Codec {
	return app.cdc
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *TerraApp) GetSubspace(moduleName string) params.Subspace {
	return app.subspaces[moduleName]
}

// SimulationManager implements the SimulationApp interface
func (app *TerraApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetTreasuryKeeper is test purpose function to return treasury keeper
func (app *TerraApp) GetTreasuryKeeper() treasury.Keeper {
	return app.treasuryKeeper
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
