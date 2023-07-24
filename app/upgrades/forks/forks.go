package forks

import (
	"fmt"

	"github.com/classic-terra/core/v2/app/keepers"
	core "github.com/classic-terra/core/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

func runForkLogicSwapDisable(ctx sdk.Context, keppers *keepers.AppKeepers, _ *module.Manager) {
	if ctx.ChainID() == core.ColumbusChainID {
		// Make min spread to 100% to disable swap
		params := keppers.MarketKeeper.GetParams(ctx)
		params.MinStabilitySpread = sdk.OneDec()
		keppers.MarketKeeper.SetParams(ctx, params)

		// Disable IBC Channels
		channelIDs := []string{
			"channel-1",  // Osmosis
			"channel-49", // Crescent
			"channel-20", // Juno
		}
		for _, channelID := range channelIDs {
			channel, found := keppers.IBCKeeper.ChannelKeeper.GetChannel(ctx, ibctransfertypes.PortID, channelID)
			if !found {
				panic(fmt.Sprintf("%s not found", channelID))
			}

			channel.State = ibcchanneltypes.CLOSED
			keppers.IBCKeeper.ChannelKeeper.SetChannel(ctx, ibctransfertypes.PortID, channelID, channel)
		}
	}
}

func runForkLogicIbcEnable(ctx sdk.Context, keppers *keepers.AppKeepers, _ *module.Manager) {
	if ctx.ChainID() == core.ColumbusChainID {
		// Enable IBC Channels
		channelIDs := []string{
			"channel-1",  // Osmosis
			"channel-49", // Crescent
			"channel-20", // Juno
		}
		for _, channelID := range channelIDs {
			channel, found := keppers.IBCKeeper.ChannelKeeper.GetChannel(ctx, ibctransfertypes.PortID, channelID)
			if !found {
				panic(fmt.Sprintf("%s not found", channelID))
			}

			channel.State = ibcchanneltypes.OPEN
			keppers.IBCKeeper.ChannelKeeper.SetChannel(ctx, ibctransfertypes.PortID, channelID, channel)
		}
	}
}

func runForkLogicVersionMapEnable(ctx sdk.Context, keppers *keepers.AppKeepers, mm *module.Manager) {
	// trigger SetModuleVersionMap in upgrade keeper at the VersionMapEnableHeight
	if ctx.ChainID() == core.ColumbusChainID {
		keppers.UpgradeKeeper.SetModuleVersionMap(ctx, mm.GetVersionMap())
	}
}
