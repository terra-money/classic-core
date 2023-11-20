package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &types.MsgCreateValidator{}, "staking/MsgCreateValidator")
	legacy.RegisterAminoMsg(cdc, &types.MsgEditValidator{}, "staking/MsgEditValidator")
	legacy.RegisterAminoMsg(cdc, &types.MsgDelegate{}, "staking/MsgDelegate")
	legacy.RegisterAminoMsg(cdc, &types.MsgUndelegate{}, "staking/MsgUndelegate")
	legacy.RegisterAminoMsg(cdc, &types.MsgBeginRedelegate{}, "staking/MsgBeginRedelegate")

	cdc.RegisterConcrete(&types.StakeAuthorization{}, "msgauth/StakeAuthorization", nil)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
