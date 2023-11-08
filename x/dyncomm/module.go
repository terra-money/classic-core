package dyncomm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/classic-terra/core/v2/x/dyncomm/client/cli"
	"github.com/classic-terra/core/v2/x/dyncomm/keeper"
	"github.com/classic-terra/core/v2/x/dyncomm/types"
	"github.com/classic-terra/core/v2/x/market/simulation"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the dyncomm module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// AppModule implements an application module for the dyncomm module.
type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	stakingKeeper types.StakingKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	stakingKeeper types.StakingKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc},
		keeper:         keeper,
		stakingKeeper:  stakingKeeper,
	}
}

// Name returns the module's name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the dyncomm
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the dyncomm module.
func (AppModuleBasic) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	return nil
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the dyncomm module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the oracle module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command for the dyncom module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// ExportGenesis returns the exported genesis state as raw bytes for the dyncomm
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// InitGenesis performs genesis initialization for the dyncomm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, &genesisState)

	return nil
}

// LegacyQuerierHandler returns the dyncomm module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewLegacyQuerier(am.keeper, legacyQuerierCdc)
}

// QuerierRoute returns the dyncomm module's querier route name.
func (AppModule) QuerierRoute() string { return types.QuerierRoute }

// RegisterInvariants performs a no-op for dyncomm module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// no msg server for this module
	querier := keeper.NewQuerier(am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

// NewHandler returns an sdk.Handler for the dyncomm module.
func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

// Route returns the message routing key for the dyncomm module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, nil)
}

// GenerateGenesisState creates a randomized GenState of the dyncomm module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	// workaround so that the staking module
	// simulation would not fail
	dyncommGenesis := types.DefaultGenesisState()
	params := types.DefaultParams()
	params.Cap = sdk.ZeroDec()
	dyncommGenesis.Params = params
	bz, err := json.MarshalIndent(&dyncommGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected default dyncomm parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(dyncommGenesis)
}

// ProposalContents returns all the dyncomm content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{}
}

// RandomizedParams creates randomized dyncomm param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	// workaround to make the sim work with staking module
	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder for dyncomm module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the dyncomm module operations with their respective weights.
func (am AppModule) WeightedOperations(module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// EndBlock returns the end blocker for the dyncomm module.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
