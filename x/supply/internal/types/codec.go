package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/terra-money/core/x/auth"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&supply.ModuleAccount{}, "supply/ModuleAccount", nil)
	cdc.RegisterConcrete(&supply.Supply{}, "supply/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()

	// Register the ModuleAccount type as a GenesisAccount so that when no
	// concrete GenesisAccount types exist and **default** genesis state is used,
	// the genesis state will serialize correctly.
	auth.RegisterAccountTypeCodec(&supply.ModuleAccount{}, "supply/ModuleAccount")
}
