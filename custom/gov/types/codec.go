package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*govv1beta1.Content)(nil), nil)
	legacy.RegisterAminoMsg(cdc, &govv1beta1.MsgSubmitProposal{}, "gov/MsgSubmitProposal")
	legacy.RegisterAminoMsg(cdc, &govv1beta1.MsgDeposit{}, "gov/MsgDeposit")
	legacy.RegisterAminoMsg(cdc, &govv1beta1.MsgVote{}, "gov/MsgVote")
	legacy.RegisterAminoMsg(cdc, &govv1beta1.MsgVoteWeighted{}, "gov/MsgVoteWeighted")
	cdc.RegisterConcrete(&govv1beta1.TextProposal{}, "gov/TextProposal", nil)
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
