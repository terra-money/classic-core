package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*govtypes.Content)(nil), nil)
	cdc.RegisterConcrete(&govtypes.MsgSubmitProposal{}, "gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(&govtypes.MsgDeposit{}, "gov/MsgDeposit", nil)
	cdc.RegisterConcrete(&govtypes.MsgVote{}, "gov/MsgVote", nil)
	cdc.RegisterConcrete(&govtypes.MsgVoteWeighted{}, "gov/MsgVoteWeighted", nil)
	cdc.RegisterConcrete(&govtypes.TextProposal{}, "gov/TextProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterProposalTypeCodec(o interface{}, name string) {
	amino.RegisterConcrete(o, name, nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
