package forks

import (
	"fmt"

	"github.com/classic-terra/core/v2/app/keepers"
	core "github.com/classic-terra/core/v2/types"
	"github.com/classic-terra/core/v2/x/market/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
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

func forkLogicFixMinCommission(ctx sdk.Context, keepers *keepers.AppKeepers, mm *module.Manager) {
	MinCommissionRate := sdk.NewDecWithPrec(5, 2)

	space, exist := keepers.ParamsKeeper.GetSubspace(stakingtypes.StoreKey)
	if !exist {
		panic("staking subspace is not found, breaking the chain anyway so panic")
	}

	if space.HasKeyTable() {
		space.Set(ctx, stakingtypes.KeyMinCommissionRate, MinCommissionRate)
	} else {
		space.WithKeyTable(types.ParamKeyTable())
		space.Set(ctx, stakingtypes.KeyMinCommissionRate, MinCommissionRate)
	}

	keepers.StakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		val := validator.(stakingtypes.Validator)
		rate := val.Commission.Rate
		maxRate := val.Commission.MaxRate

		if maxRate.LT(MinCommissionRate) {
			maxRate = MinCommissionRate
		}

		if rate.LT(MinCommissionRate) {
			rate = MinCommissionRate
		}

		val.Commission = stakingtypes.NewCommission(
			rate,
			maxRate,
			val.Commission.MaxChangeRate,
		)

		keepers.StakingKeeper.SetValidator(ctx, val)

		return false
	})
}

func runForkLogicFixMinCommission(ctx sdk.Context, keepers *keepers.AppKeepers, mm *module.Manager) {
	if ctx.ChainID() != core.ColumbusChainID {
		return
	}
	forkLogicFixMinCommission(ctx, keepers, mm)
}

func runForkLogicFixMinCommissionRebel(ctx sdk.Context, keepers *keepers.AppKeepers, mm *module.Manager) {
	if ctx.ChainID() != core.RebelChainID {
		return
	}
	forkLogicFixMinCommission(ctx, keepers, mm)
}
