package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/terra-money/core/x/gov"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(TaxRateUpdateProposal{}, "treasury/TaxRateUpdateProposal", nil)
	cdc.RegisterConcrete(RewardWeightUpdateProposal{}, "treasury/RewardWeightUpdateProposal", nil)
}

// ModuleCdc defines generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()

	gov.RegisterProposalTypeCodec(TaxRateUpdateProposal{}, "treasury/TaxRateUpdateProposal")
	gov.RegisterProposalTypeCodec(RewardWeightUpdateProposal{}, "treasury/RewardWeightUpdateProposal")
}
