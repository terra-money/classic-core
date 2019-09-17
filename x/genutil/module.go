// Package genutil has moduleCdc, we have to implement codec used part same as origin
package genutil

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	*CosmosModuleCdc = *ModuleCdc // nolint
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(GenesisState{})
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	CosmosAppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetTxCmd(cdc)
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper   AccountKeeper
	stakingKeeper   StakingKeeper
	deliverTx       func(abci.RequestDeliverTx) abci.ResponseDeliverTx
	cosmosAppModule module.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper AccountKeeper,
	stakingKeeper StakingKeeper, deliverTx func(abci.RequestDeliverTx) abci.ResponseDeliverTx) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic:  AppModuleBasic{},
		accountKeeper:   accountKeeper,
		stakingKeeper:   stakingKeeper,
		deliverTx:       deliverTx,
		cosmosAppModule: NewCosmosAppModule(accountKeeper, stakingKeeper, deliverTx),
	})
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
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, ModuleCdc, am.stakingKeeper, am.deliverTx, genesisState)
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
