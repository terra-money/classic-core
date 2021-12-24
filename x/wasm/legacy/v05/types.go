package v05

import (
	encoding_json "encoding/json"
)

const (
	// ModuleName nolint
	ModuleName = "wasm"
)

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	// Address is the address of the contract
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	// Creator is the contract creator address
	Creator string `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty" yaml:"creator"`
	// Admin is who can execute the contract migration
	Admin string `protobuf:"bytes,3,opt,name=admin,proto3" json:"admin,omitempty" yaml:"admin"`
	// CodeID is the reference to the stored Wasm code
	CodeID uint64 `protobuf:"varint,4,opt,name=code_id,json=codeId,proto3" json:"code_id,omitempty" yaml:"code_id"`
	// InitMsg is the raw message used when instantiating a contract
	InitMsg encoding_json.RawMessage `protobuf:"bytes,5,opt,name=init_msg,json=initMsg,proto3,casttype=encoding/json.RawMessage" json:"init_msg,omitempty" yaml:"init_msg"`
}
