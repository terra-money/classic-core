// nolint
package supply

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

const (
	ModuleName   = supply.ModuleName
	StoreKey     = supply.StoreKey
	RouterKey    = supply.RouterKey
	QuerierRoute = supply.QuerierRoute
	Minter       = supply.Minter
	Burner       = supply.Burner
	Staking      = supply.Staking
)

var (
	// functions aliases
	RegisterInvariants    = supply.RegisterInvariants
	AllInvariants         = supply.AllInvariants
	TotalSupply           = supply.TotalSupply
	NewKeeper             = supply.NewKeeper
	NewQuerier            = supply.NewQuerier
	SupplyKey             = supply.SupplyKey
	NewModuleAddress      = supply.NewModuleAddress
	NewEmptyModuleAccount = supply.NewEmptyModuleAccount
	NewModuleAccount      = supply.NewModuleAccount
	NewGenesisState       = supply.NewGenesisState
	DefaultGenesisState   = supply.DefaultGenesisState
	NewSupply             = supply.NewSupply
	DefaultSupply         = supply.DefaultSupply
	NewCosmosAppModule    = supply.NewAppModule

	// variable aliases
	CosmosModuleCdc = supply.ModuleCdc
)

type (
	Keeper               = supply.Keeper
	ModuleAccount        = supply.ModuleAccount
	GenesisState         = supply.GenesisState
	Supply               = supply.Supply
	CosmosAppModule      = supply.AppModule
	CosmosAppModuleBasic = supply.AppModuleBasic
	ModuleAccountI       = exported.ModuleAccountI
)
