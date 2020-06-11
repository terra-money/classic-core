package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&authtypes.BaseVestingAccount{}, "core/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&LazyGradedVestingAccount{}, "core/LazyGradedVestingAccount", nil)
}

// VestingCdc module wide codec
var VestingCdc *codec.Codec

func init() {
	VestingCdc = codec.New()
	RegisterCodec(VestingCdc)
	VestingCdc.Seal()
}
