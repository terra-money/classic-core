package bank

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/bank/client/cli"
	"github.com/terra-money/core/x/bank/client/rest"
	"github.com/terra-money/core/x/bank/internal/types"
	"github.com/terra-money/core/x/bank/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the bank module.
type AppModuleBasic struct{}

// Name returns the bank module's name
func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

// RegisterCodec registers the bank module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	*CosmosModuleCdc = *ModuleCdc // nolint
}

// DefaultGenesis returns default genesis state as raw bytes for the bank
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return CosmosAppModuleBasic{}.DefaultGenesis()
}

// ValidateGenesis performs genesis state validation for the bank module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return CosmosAppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the bank module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	rest.RegisterRoutes(cliCtx, route)
}

// GetTxCmd returns the root tx command for the bank module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the bank module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

//___________________________

// AppModule implements an application module for the bank module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
	accountKeeper   types.AccountKeeper
	keeper          Keeper
	supplyKeeper    types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper types.AccountKeeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(keeper, accountKeeper),
		accountKeeper:   accountKeeper,
		keeper:          keeper,
		supplyKeeper:    supplyKeeper,
	}
}

// Name returns the bank module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the bank module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the bank module.
func (am AppModule) Route() string {
	return am.cosmosAppModule.Route()
}

// NewHandler returns an sdk.Handler for the bank module.
func (am AppModule) NewHandler() sdk.Handler {
	return am.cosmosAppModule.NewHandler()
}

// QuerierRoute returns the bank module's querier route name.
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// NewQuerierHandler returns the bank module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier { return am.cosmosAppModule.NewQuerierHandler() }

// InitGenesis performs genesis initialization for the bank module.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// check if the burn module account exists
	burnModuleAcc := am.supplyKeeper.GetModuleAccount(ctx, types.BurnModuleName)
	if burnModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BurnModuleName))
	}

	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the bank
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// BeginBlock returns the begin blocker for the bank module.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// EndBlock returns the end blocker for the bank module.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper, am.supplyKeeper)
	return am.cosmosAppModule.EndBlock(ctx, rbb)
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the auth module
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	am.cosmosAppModule.GenerateGenesisState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []sim.WeightedProposalContent {
	return am.cosmosAppModule.ProposalContents(simState)
}

// RandomizedParams creates randomized auth param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return am.cosmosAppModule.RandomizedParams(r)
}

// RegisterStoreDecoder registers a decoder for auth module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	am.cosmosAppModule.RegisterStoreDecoder(sdr)
}

// WeightedOperations doesn't return any auth module operation.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	// We implement our own bank operations to prevent insufficient fee due to tax
	return simulation.WeightedOperations(
		simState.AppParams, simState.Cdc, am.accountKeeper, am.keeper,
	)
}
