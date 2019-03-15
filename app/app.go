package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"terra/x/budget"
	"terra/x/market"
	"terra/x/mint"
	"terra/x/oracle"
	"terra/x/pay"
	"terra/x/treasury"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
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
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
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
	keyMint          *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	paramsKeeper        params.Keeper
	oracleKeeper        oracle.Keeper
	treasuryKeeper      treasury.Keeper
	marketKeeper        market.Keeper
	budgetKeeper        budget.Keeper
	mintKeeper          mint.Keeper
}

// NewTerraApp returns a reference to an initialized TerraApp.
func NewTerraApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, baseAppOptions ...func(*bam.BaseApp)) *TerraApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	var app = &TerraApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyParams:        sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:       sdk.NewTransientStoreKey(params.TStoreKey),
		keyOracle:        sdk.NewKVStoreKey(oracle.StoreKey),
		keyTreasury:      sdk.NewKVStoreKey(treasury.StoreKey),
		keyMarket:        sdk.NewKVStoreKey(market.StoreKey),
		keyBudget:        sdk.NewKVStoreKey(budget.StoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

	// define the accountKeeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount, // target store
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)
	// add handlers
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	app.paramsKeeper = params.NewKeeper(
		app.cdc,
		app.keyParams, app.tkeyParams,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		app.bankKeeper, app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakingKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		app.keyMint,
		app.bankKeeper,
		app.accountKeeper,
	)
	app.marketKeeper = market.NewKeeper(
		app.oracleKeeper,
		app.mintKeeper,
	)
	app.treasuryKeeper = treasury.NewKeeper(
		app.cdc,
		app.keyTreasury,
		app.accountKeeper,
		app.mintKeeper,
		app.marketKeeper,
		app.paramsKeeper.Subspace(treasury.DefaultParamspace),
	)
	app.oracleKeeper = oracle.NewKeeper(
		app.cdc,
		app.keyOracle,
		stakingKeeper.GetValidatorSet(),
		app.paramsKeeper.Subspace(oracle.DefaultParamspace),
	)
	app.budgetKeeper = budget.NewKeeper(
		app.cdc,
		app.keyBudget,
		app.bankKeeper,
		stakingKeeper.GetValidatorSet(),
		app.paramsKeeper.Subspace(budget.DefaultParamspace),
	)

	// register the staking hooks
	// NOTE: The stakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	// register message routes
	app.Router().
		AddRoute(bank.RouterKey, pay.NewHandler(app.bankKeeper, app.treasuryKeeper, app.feeCollectionKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
		AddRoute(oracle.RouterKey, oracle.NewHandler(app.oracleKeeper)).
		AddRoute(budget.RouterKey, budget.NewHandler(app.budgetKeeper)).
		AddRoute(market.RouterKey, market.NewHandler(app.marketKeeper))

	app.QueryRouter().
		//AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc)).
		AddRoute(treasury.RouterKey, treasury.NewQuerier(app.treasuryKeeper)).
		AddRoute(oracle.RouterKey, oracle.NewQuerier(app.oracleKeeper)).
		AddRoute(budget.RouterKey, budget.NewQuerier(app.budgetKeeper))

	// initialize BaseApp
	app.MountStores(
		app.keyMain, app.keyAccount, app.keyStaking, app.keyDistr,
		app.keySlashing, app.keyFeeCollection, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistr, app.keyMarket,
		app.keyOracle, app.keyTreasury, app.keyBudget, app.keyMint,
	)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
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
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)
	budget.RegisterCodec(cdc)
	market.RegisterCodec(cdc)
	treasury.RegisterCodec(cdc)
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
	validatorUpdates, tags := staking.EndBlocker(ctx, app.stakingKeeper)

	_, rewardees, oracleTags := oracle.EndBlocker(ctx, app.oracleKeeper)
	tags = append(tags, oracleTags...)
	app.treasuryKeeper.ProcessClaims(ctx, treasury.OracleClaimClass, rewardees)

	claimants, budgetTags := budget.EndBlocker(ctx, app.budgetKeeper)
	tags = append(tags, budgetTags...)
	app.treasuryKeeper.ProcessClaims(ctx, treasury.BudgetClaimClass, claimants)

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
	genesisState.Sanitize()

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc = app.accountKeeper.NewAccount(ctx, acc) // set account number
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	// initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())
	treasury.InitGenesis(ctx, app.treasuryKeeper, genesisState.TreasuryData)
	budget.InitGenesis(ctx, app.budgetKeeper, genesisState.BudgetData)
	oracle.InitGenesis(ctx, app.oracleKeeper, genesisState.OracleData)

	// validate genesis state
	if err := TerraValidateGenesisState(genesisState); err != nil {
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

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
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

// LoadHeight loads a particular height
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
func (h StakingHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorCreated(ctx, valAddr)
	h.sh.AfterValidatorCreated(ctx, valAddr)
}
func (h StakingHooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.BeforeValidatorModified(ctx, valAddr)
	h.sh.BeforeValidatorModified(ctx, valAddr)
}
func (h StakingHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorRemoved(ctx, consAddr, valAddr)
	h.sh.AfterValidatorRemoved(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBonded(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBonded(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationCreated(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationCreated(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}
func (h StakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.AfterDelegationModified(ctx, delAddr, valAddr)
	h.sh.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.dh.BeforeValidatorSlashed(ctx, valAddr, fraction)
	h.sh.BeforeValidatorSlashed(ctx, valAddr, fraction)
}
