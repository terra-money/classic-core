package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "core/Account", nil)
	cdc.RegisterInterface((*auth.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseVestingAccount{}, "core/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&auth.ContinuousVestingAccount{}, "core/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&auth.DelayedVestingAccount{}, "core/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(auth.StdTx{}, "core/StdTx", nil)
	cdc.RegisterConcrete(&LazySchedule{}, "core/Schedule", nil)
	cdc.RegisterConcrete(&VestingSchedule{}, "core/VestingSchedule", nil)
	cdc.RegisterConcrete(&BaseLazyGradedVestingAccount{}, "core/LazyGradedVestingAccount", nil)
}

// ModuleCdc defines module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.RegisterConcrete(&supply.ModuleAccount{}, "supply/ModuleAccount", nil)
	ModuleCdc.Seal()
}
