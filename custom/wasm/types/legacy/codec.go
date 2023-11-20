package legacy

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the wasm types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgStoreCode{}, "wasm/MsgStoreCode")
	legacy.RegisterAminoMsg(cdc, &MsgMigrateCode{}, "wasm/MsgMigrateCode")
	legacy.RegisterAminoMsg(cdc, &MsgInstantiateContract{}, "wasm/MsgInstantiateContract")
	legacy.RegisterAminoMsg(cdc, &MsgExecuteContract{}, "wasm/MsgExecuteContract")
	legacy.RegisterAminoMsg(cdc, &MsgMigrateContract{}, "wasm/MsgMigrateContract")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateContractAdmin{}, "wasm/MsgUpdateContractAdmin")
	legacy.RegisterAminoMsg(cdc, &MsgClearContractAdmin{}, "wasm/MsgClearContractAdmin")
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
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
