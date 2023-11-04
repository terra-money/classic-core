package v6

import (
	"github.com/classic-terra/core/v2/app/upgrades"
	dyncommtypes "github.com/classic-terra/core/v2/x/dyncomm/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const UpgradeName = "v6"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			dyncommtypes.StoreKey,
		},
	},
}
