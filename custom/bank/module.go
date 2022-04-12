package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	customsim "github.com/terra-money/core/custom/bank/simulation"
	customtypes "github.com/terra-money/core/custom/bank/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the distribution module.
type AppModuleBasic struct {
	bank.AppModuleBasic
}

// RegisterLegacyAminoCodec registers the bank module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	customtypes.RegisterLegacyAminoCodec(cdc)
	*types.ModuleCdc = *customtypes.ModuleCdc // nolint
}

// AppModule implements an application module for the bank module.
type AppModule struct {
	bank.AppModule
	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModule:     bank.NewAppModule(cdc, keeper, accountKeeper),
		keeper:        keeper,
		accountKeeper: accountKeeper,
	}
}

// GenerateGenesisState creates a randomized GenState of the bank module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	customsim.RandomizedGenState(simState)
}

// WeightedOperations return random bank module operation.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// We implement our own bank operations to prevent insufficient fee due to tax
	return customsim.WeightedOperations(
		simState.AppParams, simState.Cdc, am.accountKeeper, am.keeper,
	)
}
