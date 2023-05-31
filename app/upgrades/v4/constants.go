package v3

import (
	"github.com/classic-terra/core/app/upgrades"
	feesharetypes "github.com/classic-terra/core/x/feeshare/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
)

const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{Added: []string{feesharetypes.StoreKey, icahosttypes.StoreKey}},
}
