package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/terra-money/core/custom/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&types.MsgWithdrawDelegatorReward{}, "distribution/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(&types.MsgWithdrawValidatorCommission{}, "distribution/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(&types.MsgSetWithdrawAddress{}, "distribution/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(&types.MsgFundCommunityPool{}, "distribution/MsgFundCommunityPool", nil)
	cdc.RegisterConcrete(&types.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/distribution module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/distribution and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()

	govtypes.RegisterProposalTypeCodec(types.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal")
}
