// nolint
package genutil

import (
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

const (
	ModuleName = genutil.ModuleName
)

var (
	// functions aliases
	NewCosmosAppModule           = genutil.NewAppModule
	InitGenesis                  = genutil.InitGenesis
	InitializeNodeValidatorFiles = genutil.InitializeNodeValidatorFiles
	NewInitConfig                = genutil.NewInitConfig
	ExportGenesisFileWithTime    = genutil.ExportGenesisFileWithTime
	GenesisStateFromGenFile      = genutil.GenesisStateFromGenFile
	GenAppStateFromConfig        = genutil.GenAppStateFromConfig
	ExportGenesisFile            = genutil.ExportGenesisFile
	CosmosModuleCdc              = genutil.ModuleCdc
)

type (
	GenesisState         = genutil.GenesisState
	CosmosAppModule      = genutil.AppModule
	CosmosAppModuleBasic = genutil.AppModuleBasic
)
