// nolint
package upgrade

import (
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade"
)

const (
	ModuleName                        = upgrade.ModuleName
	RouterKey                         = upgrade.RouterKey
	StoreKey                          = upgrade.StoreKey
	QuerierKey                        = upgrade.QuerierKey
	PlanByte                          = upgrade.PlanByte
	DoneByte                          = upgrade.DoneByte
	ProposalTypeSoftwareUpgrade       = upgrade.ProposalTypeSoftwareUpgrade
	ProposalTypeCancelSoftwareUpgrade = upgrade.ProposalTypeCancelSoftwareUpgrade
	QueryCurrent                      = upgrade.QueryCurrent
	QueryApplied                      = upgrade.QueryApplied
)

var (
	// functions aliases
	PlanKey                           = upgrade.PlanKey
	NewSoftwareUpgradeProposal        = upgrade.NewSoftwareUpgradeProposal
	NewCancelSoftwareUpgradeProposal  = upgrade.NewCancelSoftwareUpgradeProposal
	NewQueryAppliedParams             = upgrade.NewQueryAppliedParams
	NewKeeper                         = upgrade.NewKeeper
	NewQuerier                        = upgrade.NewQuerier
	NewSoftwareUpgradeProposalHandler = upgrade.NewSoftwareUpgradeProposalHandler

	NewCosmosAppModule = upgrade.NewAppModule
)

type (
	UpgradeHandler                = upgrade.UpgradeHandler
	Plan                          = upgrade.Plan
	SoftwareUpgradeProposal       = upgrade.SoftwareUpgradeProposal
	CancelSoftwareUpgradeProposal = upgrade.CancelSoftwareUpgradeProposal
	QueryAppliedParams            = upgrade.QueryAppliedParams
	Keeper                        = upgrade.Keeper

	CosmosAppModuleBasic = upgrade.AppModuleBasic
	CosmosAppModule      = upgrade.AppModule
)
