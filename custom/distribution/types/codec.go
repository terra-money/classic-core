package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/classic-terra/core/v2/custom/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// distribution/MsgWithdrawDelegationReward and distribution/MsgWithdrawValidatorCommission
	// will not be supported by Ledger signing due to the overflow of length of the name
	cdc.RegisterConcrete(&types.MsgWithdrawDelegatorReward{}, "distribution/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(&types.MsgWithdrawValidatorCommission{}, "distribution/MsgWithdrawValidatorCommission", nil)
	legacy.RegisterAminoMsg(cdc, &types.MsgSetWithdrawAddress{}, "distribution/MsgModifyWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &types.MsgFundCommunityPool{}, "distribution/MsgFundCommunityPool")
	cdc.RegisterConcrete(&types.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal", nil)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()

	govtypes.RegisterProposalTypeCodec(types.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal")
}
