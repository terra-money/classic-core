package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"terra/types/tax"
	"terra/x/budget"
	"terra/x/market"
	"terra/x/oracle"
	"terra/x/treasury"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	appName = "TerraApp"
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.terracli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.terrad")
)

// Extended ABCI application
type TerraApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyBank          *sdk.KVStoreKey
	keyStake         *sdk.KVStoreKey
	tkeyStake        *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey
	keyOracle        *sdk.KVStoreKey
	keyTreasury      *sdk.KVStoreKey
	keyMarket        *sdk.KVStoreKey
	keyBudget        *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          tax.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper
	oracleKeeper        oracle.Keeper
	treasuryKeeper      treasury.Keeper
	marketKeeper        market.Keeper
	budgetKeeper        budget.Keeper
}

// NewTerraApp returns a reference to an initialized TerraApp.
func NewTerraApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, baseAppOptions ...func(*bam.BaseApp)) *TerraApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	var app = &TerraApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyBank:          sdk.NewKVStoreKey("bank"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		tkeyStake:        sdk.NewTransientStoreKey("transient_stake"),
		keyDistr:         sdk.NewKVStoreKey("distr"),
		tkeyDistr:        sdk.NewTransientStoreKey("transient_distr"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),
		keyOracle:        sdk.NewKVStoreKey("oracle"),
		keyTreasury:      sdk.NewKVStoreKey("treasury"),
		keyMarket:        sdk.NewKVStoreKey("market"),
		keyBudget:        sdk.NewKVStoreKey("budget"),
	}

	// define the accountKeeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)
	// add handlers
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	app.bankKeeper = tax.NewBaseKeeper(
		app.keyBank,
		app.cdc,
		app.accountKeeper,
		app.feeCollectionKeeper,
	)
	app.paramsKeeper = params.NewKeeper(
		app.cdc,
		app.keyParams, app.tkeyParams,
	)
	stakeKeeper := stake.NewKeeper(
		app.cdc,
		app.keyStake, app.tkeyStake,
		app.bankKeeper, app.paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakeKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.treasuryKeeper = treasury.NewKeeper(
		app.keyTreasury,
		app.cdc,
		app.bankKeeper,
	)
	app.oracleKeeper = oracle.NewKeeper(
		app.keyOracle,
		cdc,
		app.treasuryKeeper,
		stakeKeeper.GetValidatorSet(),
		app.paramsKeeper.Subspace(oracle.DefaultParamspace),
	)
	app.marketKeeper = market.NewKeeper(
		app.oracleKeeper,
		app.treasuryKeeper,
		app.bankKeeper,
	)
	app.budgetKeeper = budget.NewKeeper(
		app.keyBudget,
		app.cdc, app.bankKeeper,
		app.treasuryKeeper,
		stake.DefaultCodespace,
		stakeKeeper.GetValidatorSet(),
		app.paramsKeeper.Subspace(budget.DefaultParamspace),
	)

	// register the staking hooks
	// NOTE: The stakeKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakeKeeper = *stakeKeeper.SetHooks(
		NewStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("distr", distr.NewHandler(app.distrKeeper)).
		AddRoute("slashing", slashing.NewHandler(app.slashingKeeper)).
		AddRoute("oracle", oracle.NewHandler(app.oracleKeeper)).
		AddRoute("budget", budget.NewHandler(app.budgetKeeper)).
		AddRoute("market", market.NewHandler(app.marketKeeper))

	app.QueryRouter().
		AddRoute("budget", budget.NewQuerier(app.budgetKeeper)).
		AddRoute("stake", stake.NewQuerier(app.stakeKeeper, app.cdc))

	// initialize BaseApp
	app.MountStores(app.keyMain, app.keyAccount, app.keyStake, app.keyDistr, app.keyBank,
		app.keySlashing, app.keyFeeCollection, app.keyParams, app.keyMarket, app.keyOracle, app.keyTreasury, app.keyBudget)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
	app.MountStoresTransient(app.tkeyParams, app.tkeyStake, app.tkeyDistr)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// application updates every end block
func (app *TerraApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {

	// distribute rewards for the previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// slash anyone who double signed.
	// NOTE: This should happen after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool,
	// so as to keep the CanWithdrawInvariant invariant.
	// TODO: This should really happen at EndBlocker.
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
// nolint: unparam
func (app *TerraApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

	validatorUpdates, tags := stake.EndBlocker(ctx, app.stakeKeeper)

	oracleTags := oracle.EndBlocker(ctx, app.oracleKeeper)
	tags = append(tags, oracleTags...)

	budgetTags := budget.EndBlocker(ctx, app.budgetKeeper)
	tags = append(tags, budgetTags...)

	treasuryTags := treasury.EndBlocker(ctx, app.treasuryKeeper)
	tags = append(tags, treasuryTags...)

	// TODO: request fixing it to comsmos guys
	//app.assertRuntimeInvariants()

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// initialize store from a genesis state
func (app *TerraApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	// sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})
	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// load the initial stake information
	validators, err := stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	// initialize module-specific stores
	auth.InitGenesis(ctx, app.feeCollectionKeeper, genesisState.AuthData)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	treasury.InitGenesis(ctx, app.treasuryKeeper, genesisState.TreasuryData)
	budget.InitGenesis(ctx, app.budgetKeeper, genesisState.BudgetData)
	oracle.InitGenesis(ctx, app.oracleKeeper, genesisState.OracleData)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	// validate genesis state
	err = TerraValidateGenesisState(genesisState)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return validators
}

// custom logic for Terra initialization
func (app *TerraApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	validators := app.initFromGenesisState(ctx, genesisState)

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	// assert runtime invariants
	//app.assertRuntimeInvariants()

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

// load a particular height
func (app *TerraApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

//______________________________________________________________________________________________

var _ sdk.StakingHooks = StakingHooks{}

// StakingHooks contains combined distribution and slashing hooks needed for the
// staking module.
type StakingHooks struct {
	dh distr.Hooks
	sh slashing.Hooks
}

func NewStakingHooks(dh distr.Hooks, sh slashing.Hooks) StakingHooks {
	return StakingHooks{dh, sh}
}

// nolint
func (h StakingHooks) OnValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.OnValidatorCreated(ctx, valAddr)
	h.sh.OnValidatorCreated(ctx, valAddr)
}
func (h StakingHooks) OnValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.OnValidatorModified(ctx, valAddr)
	h.sh.OnValidatorModified(ctx, valAddr)
}
func (h StakingHooks) OnValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorRemoved(ctx, consAddr, valAddr)
	h.sh.OnValidatorRemoved(ctx, consAddr, valAddr)
}
func (h StakingHooks) OnValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorBonded(ctx, consAddr, valAddr)
	h.sh.OnValidatorBonded(ctx, consAddr, valAddr)
}
func (h StakingHooks) OnValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
	h.sh.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
}
func (h StakingHooks) OnValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.sh.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (h StakingHooks) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationCreated(ctx, delAddr, valAddr)
	h.sh.OnDelegationCreated(ctx, delAddr, valAddr)
}
func (h StakingHooks) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationSharesModified(ctx, delAddr, valAddr)
	h.sh.OnDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationRemoved(ctx, delAddr, valAddr)
	h.sh.OnDelegationRemoved(ctx, delAddr, valAddr)
}
