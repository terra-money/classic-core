package keepers

import (
	"github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"

	"github.com/classic-terra/core/x/treasury"
	treasurytypes "github.com/classic-terra/core/x/treasury/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (appKeepers *AppKeepers) getGovRouter() govtypes.Router {
	govRouter := govtypes.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(treasurytypes.RouterKey, treasury.NewProposalHandler(appKeepers.TreasuryKeeper))

	return govRouter
}

func (appKeepers *AppKeepers) setIBCRouter() {
	transferModule := transfer.NewAppModule(appKeepers.TransferKeeper)

	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)
}
