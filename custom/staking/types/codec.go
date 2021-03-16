package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&types.MsgCreateValidator{}, "staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(&types.MsgEditValidator{}, "staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(&types.MsgDelegate{}, "staking/MsgDelegate", nil)
	cdc.RegisterConcrete(&types.MsgUndelegate{}, "staking/MsgUndelegate", nil)
	cdc.RegisterConcrete(&types.MsgBeginRedelegate{}, "staking/MsgBeginRedelegate", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/staking module codec. Note, the codec should
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
	amino.Seal()
}
