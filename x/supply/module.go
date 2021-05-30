package supply

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/supply/internal/types"
	"github.com/terra-money/core/x/supply/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the supply module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the supply module's name
func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

// RegisterCodec registers the supply module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	*CosmosModuleCdc = *ModuleCdc // nolint
}

// DefaultGenesis returns default genesis state as raw bytes for the supply
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return CosmosAppModuleBasic{}.DefaultGenesis()
}

// ValidateGenesis performs genesis state validation for the supply module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return CosmosAppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the supply module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	CosmosAppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// GetTxCmd returns the root tx command for the supply module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the supply module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

//___________________________

// AppModule implements an application module for the supply module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(keeper, ak),
	}
}

// Name returns the supply module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the supply module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the supply module.
func (am AppModule) Route() string {
	return am.cosmosAppModule.Route()
}

// NewHandler returns an sdk.Handler for the supply module.
func (am AppModule) NewHandler() sdk.Handler {
	return am.cosmosAppModule.NewHandler()
}

// QuerierRoute returns the supply module's querier route name.
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// NewQuerierHandler returns the supply module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier { return am.cosmosAppModule.NewQuerierHandler() }

// InitGenesis performs genesis initialization for the supply module.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the supply
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// BeginBlock returns the begin blocker for the supply module.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// EndBlock returns the end blocker for the supply module.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, rbb)
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the auth module
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
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
	return am.cosmosAppModule.WeightedOperations(simState)
}
