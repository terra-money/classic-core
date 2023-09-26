package fork

// NOTE: Keep soft-fork heights for the history
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
	// Revert Min Commission slip during v2.2.0 upgrade
	FixMinCommissionHeight      = int64(14_890_000)
	FixMinCommissionHeightRebel = int64(16_300_000)
)
