package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// RegisterLegacyAminoCodec registers the account interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*types.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*types.AccountI)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAccount{}, "core/Account", nil)
	cdc.RegisterConcrete(&types.ModuleAccount{}, "core/ModuleAccount", nil)
	cdc.RegisterConcrete(legacytx.StdTx{}, "core/StdTx", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/oracle module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
