// nolint
package mint

import (
	"github.com/cosmos/cosmos-sdk/x/mint"
)

const (
	ModuleName        = mint.ModuleName
	StoreKey          = mint.StoreKey
	QuerierRoute      = mint.QuerierRoute
	QueryParameters   = mint.QueryParameters
	DefaultParamspace = mint.DefaultParamspace
)

var (
	// functions aliases
	NewGenesisState     = mint.NewGenesisState
	DefaultGenesisState = mint.DefaultGenesisState
	ValidateGenesis     = mint.ValidateGenesis
	ParamKeyTable       = mint.ParamKeyTable
	NewParams           = mint.NewParams
	DefaultParams       = mint.DefaultParams
	NewCosmosAppModule  = mint.NewAppModule
	NewKeeper           = mint.NewKeeper

	// variable aliases
	CosmosModuleCdc = mint.ModuleCdc
)

type (
	GenesisState         = mint.GenesisState
	Params               = mint.Params
	Keeper               = mint.Keeper
	CosmosAppModule      = mint.AppModule
	CosmosAppModuleBasic = mint.AppModuleBasic
)
