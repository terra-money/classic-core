package gov

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gov/client"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/gov/internal/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct {
	CosmosAppModuleBasic
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic(proposalHandlers ...client.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		CosmosAppModuleBasic: NewCosmosAppModuleBasic(proposalHandlers...),
	}
}

var _ module.AppModuleBasic = AppModuleBasic{}

// module name
func (am AppModuleBasic) Name() string {
	return am.CosmosAppModuleBasic.Name()
}

// register module codec
func (am AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	*CosmosModuleCdc = *ModuleCdc // nolint
}

// default genesis state
func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	// customize to set default genesis state deposit denom to uluna
	defaultGenesisState := DefaultGenesisState()
	defaultGenesisState.DepositParams.MinDeposit[0].Denom = core.MicroLunaDenom

	return ModuleCdc.MustMarshalJSON(defaultGenesisState)
}

// module validate genesis
func (am AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return am.CosmosAppModuleBasic.ValidateGenesis(bz)
}

// register rest routes
func (am AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	am.CosmosAppModuleBasic.RegisterRESTRoutes(cliCtx, route)
}

// get the root tx command of this module
func (am AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return am.CosmosAppModuleBasic.GetTxCmd(cdc)
}

// get the root query command of this module
func (am AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return am.CosmosAppModuleBasic.GetQueryCmd(cdc)
}

//___________________________
// app module for gov
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(keeper, supplyKeeper),
	}
}

// module name
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// register invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// module querier route name
func (am AppModule) Route() string {
	return am.cosmosAppModule.Route()
}

// module handler
func (am AppModule) NewHandler() sdk.Handler {
	return am.cosmosAppModule.NewHandler()
}

// module querier route name
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier { return am.cosmosAppModule.NewQuerierHandler() }

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// module end-block
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, rbb)
}
