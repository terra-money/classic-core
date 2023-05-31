package wasm

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/CosmWasm/wasmd/x/wasm"

	customcli "github.com/classic-terra/core/v2/custom/wasm/client/cli"
)

var (
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
