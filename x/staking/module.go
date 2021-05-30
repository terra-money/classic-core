package staking

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/staking/internal/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the staking module.
type AppModuleBasic struct{}

// Name returns the staking module's name
func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	*CosmosModuleCdc = *ModuleCdc // nolint
}

// DefaultGenesis returns default genesis state as raw bytes for the staking
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	// customize to set default genesis state bond denom to uluna
	defaultGenesisState := DefaultGenesisState()
	defaultGenesisState.Params.BondDenom = core.MicroLunaDenom

	return ModuleCdc.MustMarshalJSON(defaultGenesisState)
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return CosmosAppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the staking module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	CosmosAppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// GetTxCmd returns the root tx command for the staking module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the staking module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

//_____________________________________
// extra helpers

// CreateValidatorMsgHelpers - used for gen-tx
func (AppModuleBasic) CreateValidatorMsgHelpers(ipDefault string) (
	fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {
	return CosmosAppModuleBasic{}.CreateValidatorMsgHelpers(ipDefault)
}

// PrepareFlagsForTxCreateValidator - used for gen-tx
func (AppModuleBasic) PrepareFlagsForTxCreateValidator(config *cfg.Config, nodeID,
	chainID string, valPubKey crypto.PubKey) {
	CosmosAppModuleBasic{}.PrepareFlagsForTxCreateValidator(config, nodeID, chainID, valPubKey)
}

// BuildCreateValidatorMsg - used for gen-tx
func (AppModuleBasic) BuildCreateValidatorMsg(cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder) (authtypes.TxBuilder, sdk.Msg, error) {
	return CosmosAppModuleBasic{}.BuildCreateValidatorMsg(cliCtx, txBldr)
}

//___________________________

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accKeeper types.AccountKeeper,
	supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(keeper, accKeeper, supplyKeeper),
	}
}

// Name returns the staking module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() string {
	return am.cosmosAppModule.Route()
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return am.cosmosAppModule.NewHandler()
}

// QuerierRoute returns the staking module's querier route name.
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// NewQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier { return am.cosmosAppModule.NewQuerierHandler() }

// InitGenesis performs genesis initialization for the staking module. It returns
// validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// EndBlock returns the end blocker for the staking module.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
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
	return am.cosmosAppModule.WeightedOperations(simState)
}
