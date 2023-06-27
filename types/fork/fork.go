package fork

import (
	"github.com/classic-terra/core/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ColumbusLunaSwapFeeHeight = int64(5_100_000)
	BombayLunaSwapFeeHeight   = int64(6_200_000)
	ColumbusOracleFixHeight   = int64(5_701_000)
	BombayOracleFixHeight     = int64(7_000_000)
	// v0.5.20
	// SwapDisableHeight - make min spread to 100% to disable swap
	SwapDisableHeight = int64(7_607_790)
	// v0.5.21
	// BurnTaxUpgradeHeight is when taxes are allowed to go into effect
	// This will still need a parameter change proposal, but can be activated
	// anytime after this height
	BurnTaxUpgradeHeight = int64(9_346_889)
	// v0.5.23
	// IbcEnableHeight - renable IBC only, block height is approximately December 5th, 2022
	IbcEnableHeight = int64(10_542_500)
	// v1.0.5
	// VersionMapEnableHeight - set the version map to enable software upgrades, approximately February 14, 2023
	VersionMapEnableHeight = int64(11_543_150)
)

func IsBeforeLunaSwapFeeHeight(ctx sdk.Context) bool {
	return (ctx.ChainID() == types.ColumbusChainID && ctx.BlockHeight() < ColumbusLunaSwapFeeHeight) ||
		(ctx.ChainID() == types.BombayChainID && ctx.BlockHeight() < BombayLunaSwapFeeHeight)
}

func IsBeforeOracleFixHeight(ctx sdk.Context) bool {
	return (ctx.ChainID() == types.ColumbusChainID && ctx.BlockHeight() < ColumbusOracleFixHeight) ||
		(ctx.ChainID() == types.BombayChainID && ctx.BlockHeight() < BombayOracleFixHeight)
}

func IsBeforeBurnTaxUpgradeHeight(ctx sdk.Context) bool {
	return (ctx.ChainID() == types.ColumbusChainID && ctx.BlockHeight() < BurnTaxUpgradeHeight)
}
