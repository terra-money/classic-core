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
	NewGenesisState              = genutil.NewGenesisState
	NewGenesisStateFromStdTx     = genutil.NewGenesisStateFromStdTx
	NewInitConfig                = genutil.NewInitConfig
	GetGenesisStateFromAppState  = genutil.GetGenesisStateFromAppState
	SetGenesisStateInAppState    = genutil.SetGenesisStateInAppState
	GenesisStateFromGenDoc       = genutil.GenesisStateFromGenDoc
	GenesisStateFromGenFile      = genutil.GenesisStateFromGenFile
	ValidateGenesis              = genutil.ValidateGenesis
	ExportGenesisFile            = genutil.ExportGenesisFile
	CosmosModuleCdc              = genutil.ModuleCdc
)

type (
	GenesisState      = genutil.GenesisState
	AppMap            = genutil.AppMap
	MigrationCallback = genutil.MigrationCallback
	MigrationMap      = genutil.MigrationMap
	InitConfig        = genutil.InitConfig

	CosmosAppModule      = genutil.AppModule
	CosmosAppModuleBasic = genutil.AppModuleBasic
)
