package v038

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v038upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/legacy/v038"
)

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(v038upgrade.Plan{}, "upgrade/Plan", nil)
	cdc.RegisterConcrete(v038upgrade.SoftwareUpgradeProposal{}, "upgrade/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(v038upgrade.CancelSoftwareUpgradeProposal{}, "upgrade/CancelSoftwareUpgradeProposal", nil)
}
