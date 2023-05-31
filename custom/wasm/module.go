package wasm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/simulation"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	customcli "github.com/classic-terra/core/custom/wasm/client/cli"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the wasm module.
type AppModuleBasic struct {
	wasm.AppModuleBasic
}

// GetTxCmd returns the root tx command for the wasm module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return customcli.GetTxCmd()
}

type AppModule struct {
	wasm.AppModule
	appModuleBasic     AppModuleBasic
	cdc                codec.Codec
	keeper             *wasm.Keeper
	validatorSetSource keeper.ValidatorSetSource
	accountKeeper      types.AccountKeeper // for simulation
	bankKeeper         simulation.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *wasm.Keeper,
	validatorSetSource keeper.ValidatorSetSource,
	ak types.AccountKeeper,
	bk simulation.BankKeeper,
) AppModule {
	return AppModule{
		appModuleBasic:     AppModuleBasic{},
		cdc:                cdc,
		keeper:             keeper,
		validatorSetSource: validatorSetSource,
		accountKeeper:      ak,
		bankKeeper:         bk,
	}
}
