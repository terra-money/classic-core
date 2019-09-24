package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/terra-project/core/x/gov"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(distribution.MsgWithdrawDelegatorReward{}, "distribution/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(distribution.MsgWithdrawValidatorCommission{}, "distribution/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(distribution.MsgSetWithdrawAddress{}, "distribution/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(distribution.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal", nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()

	gov.RegisterProposalTypeCodec(distribution.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal")
}
