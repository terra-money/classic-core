package authz

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"

	customtypes "github.com/terra-money/core/custom/authz/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the distribution module.
type AppModuleBasic struct {
	authz.AppModuleBasic
}

// RegisterLegacyAminoCodec registers the bank module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	customtypes.RegisterLegacyAminoCodec(cdc)
}
