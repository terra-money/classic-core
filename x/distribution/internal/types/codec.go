package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/terra-money/core/x/gov"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// NOTE: type and message name inconsistency
	cdc.RegisterConcrete(distrtypes.MsgWithdrawDelegatorReward{}, "distribution/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(distrtypes.MsgWithdrawValidatorCommission{}, "distribution/MsgWithdrawValidatorCommission", nil)
	// NOTE: type and message name inconsistency
	cdc.RegisterConcrete(distrtypes.MsgSetWithdrawAddress{}, "distribution/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(distrtypes.MsgFundCommunityPool{}, "distribution/MsgFundCommunityPool", nil)
	cdc.RegisterConcrete(distrtypes.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal", nil)
}

// ModuleCdc is generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()

	gov.RegisterProposalTypeCodec(distrtypes.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal")
}
