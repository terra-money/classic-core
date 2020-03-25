package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the wasm types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgStoreCode{}, "wasm/StoreCode", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "wasm/InstantiateContract", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "wasm/ExecuteContract", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
