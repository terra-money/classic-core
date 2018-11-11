package app

import (
	"encoding/json"
	"os"

	"terra/types"

	"terra/x/auth"
	"terra/x/bank"
	"terra/x/ibc"
	"terra/x/oracle"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "TerraApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.terracli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.terrad")
)

// TerraApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type TerraApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyIBC           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyOracle        *sdk.KVStoreKey

	// manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	oracleKeeper        oracle.Keeper
}

// NewTerraApp returns a reference to a new TerraApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewTerraApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *TerraApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &TerraApp{
		cdc:              cdc,
		BaseApp:          bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyIBC:           sdk.NewKVStoreKey("ibc"),
	}

	// define and attach the mappers and keepers
	app.accountKeeper = auth.NewAccountKeeper(
		cdc,
		app.keyAccount, // target store
		func() auth.Account {
			return &types.AppAccount{}
		},
	)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, app.keyFeeCollection)
	app.bankKeeper = bank.NewBaseTaxKeeper(app.accountKeeper, app.feeCollectionKeeper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("oracle", oracle.NewHandler(app.oracleKeeper))

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

// MakeCodec creates a new codec codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)

	// register custom type
	cdc.RegisterConcrete(&types.AppAccount{}, "terra/Account", nil)

	cdc.Seal()

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *TerraApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *TerraApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	tags := oracle.EndBlocker(ctx, app.oracleKeeper)
	//validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)

	return abci.ResponseEndBlock{
		//ValidatorUpdates: validatorUpdates,
		Tags: tags,
	}
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *TerraApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *TerraApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*types.GenesisAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := types.GenesisState{Accounts: accounts}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}
