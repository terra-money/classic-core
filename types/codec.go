package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec register type codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&Schedule{}, "core/Schedule", nil)
	cdc.RegisterConcrete(&VestingSchedule{}, "core/VestingSchedule", nil)
	cdc.RegisterConcrete(&GradedVestingAccount{}, "core/GradedVestingAccount", nil)
}
