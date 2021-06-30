package app

import (
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/terra-money/core/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// authz module use this codec to get signbytes.
	// authz MsgExec can execute all message types,
	// so legacy.Cdc need to register all amino messages to get proper signature
	ModuleBasics.RegisterLegacyAminoCodec(legacy.Cdc)

	return encodingConfig
}
