package wasm

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/wasm"

	customcli "github.com/classic-terra/core/v2/custom/wasm/client/cli"
	customrest "github.com/classic-terra/core/v2/custom/wasm/client/rest"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the wasm module.
type AppModuleBasic struct {
	wasm.AppModuleBasic
}

// RegisterRESTRoutes registers the REST routes for the wasm module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx client.Context, rtr *mux.Router) {
	customrest.RegisterRoutes(cliCtx, rtr)
}

// GetTxCmd returns the root tx command for the wasm module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return customcli.GetTxCmd()
}
