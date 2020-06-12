package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/evidence/exported"
)

// ModuleCdc defines the evidence module's codec. The codec is not sealed as to
// allow other modules to register their concrete Evidence types.
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the
// evidence module.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Evidence)(nil), nil)
	cdc.RegisterConcrete(evidence.MsgSubmitEvidence{}, "evidence/MsgSubmitEvidence", nil)
	cdc.RegisterConcrete(evidence.Equivocation{}, "evidence/Equivocation", nil)
}

// RegisterEvidenceTypeCodec registers an external concrete Evidence type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitEvidence
// to be correctly Amino encoded and decoded.
func RegisterEvidenceTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
