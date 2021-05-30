package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/gov"

	msgauthexported "github.com/terra-money/core/x/msgauth/exported"
)

// ModuleCdc defines module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*gov.Content)(nil), nil)

	cdc.RegisterConcrete(gov.MsgSubmitProposal{}, "gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(gov.MsgDeposit{}, "gov/MsgDeposit", nil)
	cdc.RegisterConcrete(gov.MsgVote{}, "gov/MsgVote", nil)

	cdc.RegisterConcrete(gov.TextProposal{}, "gov/TextProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)

	msgauthexported.RegisterMsgAuthTypeCodec(gov.MsgVote{}, "gov/MsgVote")
}
