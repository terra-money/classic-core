package legacy

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the wasm types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStoreCode{}, "wasm/MsgStoreCode", nil)
	cdc.RegisterConcrete(&MsgMigrateCode{}, "wasm/MsgMigrateCode", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)
	cdc.RegisterConcrete(&MsgMigrateContract{}, "wasm/MsgMigrateContract", nil)
	cdc.RegisterConcrete(&MsgUpdateContractAdmin{}, "wasm/MsgUpdateContractAdmin", nil)
	cdc.RegisterConcrete(&MsgClearContractAdmin{}, "wasm/MsgClearContractAdmin", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// register legacy wasm msgs (to be used in archives)
	registry.RegisterInterface("terra.wasm.v1beta1.MsgInstantiateContract", (*sdk.Msg)(nil), &MsgInstantiateContract{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgExecuteContract", (*sdk.Msg)(nil), &MsgExecuteContract{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgStoreCode", (*sdk.Msg)(nil), &MsgStoreCode{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgExecuteContract", (*sdk.Msg)(nil), &MsgMigrateCode{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgMigrateCode", (*sdk.Msg)(nil), &MsgMigrateContract{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgUpdateContractAdmin", (*sdk.Msg)(nil), &MsgUpdateContractAdmin{})
	registry.RegisterInterface("terra.wasm.v1beta1.MsgClearContractAdmin", (*sdk.Msg)(nil), &MsgClearContractAdmin{})
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/market module codec. Note, the codec should
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
