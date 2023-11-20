package v3

import (
	"github.com/classic-terra/core/v2/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const UpgradeName = "v3"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV3UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
