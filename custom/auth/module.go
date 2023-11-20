package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	customsim "github.com/classic-terra/core/v2/custom/auth/simulation"
	customtypes "github.com/classic-terra/core/v2/custom/auth/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the auth module.
type AppModuleBasic struct {
	auth.AppModuleBasic
}

// RegisterLegacyAminoCodec registers the auth module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	customtypes.RegisterLegacyAminoCodec(cdc)
	*types.ModuleCdc = *customtypes.ModuleCdc
}

// ____________________________________________________________________________

// AppModule implements an application module for the auth module.
type AppModule struct {
	auth.AppModule
	accountKeeper     keeper.AccountKeeper
	randGenAccountsFn types.RandomGenesisAccountsFn
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, accountKeeper keeper.AccountKeeper, randGenAccountsFn types.RandomGenesisAccountsFn) AppModule {
	return AppModule{
		AppModule:         auth.NewAppModule(cdc, accountKeeper, randGenAccountsFn),
		accountKeeper:     accountKeeper,
		randGenAccountsFn: randGenAccountsFn,
	}
}

// GenerateGenesisState creates a randomized GenState of the auth module
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	customsim.RandomizedGenState(simState, am.randGenAccountsFn)
}
