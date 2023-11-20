package v5

import (
	"github.com/classic-terra/core/v2/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
)

const UpgradeName = "v5"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV5UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ibcfeetypes.StoreKey,
			icacontrollertypes.StoreKey,
		},
	},
}
