# Classic Terra Upgrades

This folder contains sub-folders for every classic terra upgrade. (Both state
migrations, and hard forks) It also defines upgrade & hard fork structs,
that each upgrade implements. These then get included in the application
app.go to run the upgrade.

## Upgrade History
### v1 (core v1.0.x)
__NOTE__: initial version without upgrade handler
* (soft-fork) v0.5.20 - Disable swap and IBC [#760](https://github.com/terra-money/classic-core/pull/760)
* (soft-fork) v0.5.21 - Burn Tax and allow Luna to be taxed
* (soft-fork) v0.5.22 - Limit validator power to 20% [cosmos-sdk/pull/104](https://github.com/terra-money/cosmos-sdk/pull/104)
* (soft-fork) v0.5.23 - Enable IBC
* (soft-fork) v1.0.5 - Set module version map for enabling software upgrades
### v2 (core v1.1.x)
* (treasury) Tax exemption list, burn tax split
### v3 (core v2.0.x)
* (treasury) Minimum initial deposit for governance proposals
### v4 (core v2.1.x)
* (staking) Minimum commision rate

## Upgrade types

There are two upgrade types exposed, `Upgrade` and `Fork`. An `Upgrade`
defines an upgrade that is to be acted upon by state migrations from the
SDK `x/upgrade` module. A `Fork` defines a hard fork that changes some
logic at a block height. If the goal is to have a new binary be
compatible with the old binary prior to the upgrade height, as is the
case for all classic terra `Fork`s, then all logic changes must be
height-gated or in the `BeginForkLogic` code.

```go
type Upgrade struct {
 // Upgrade version name, for the upgrade handler, e.g. `v7`
 UpgradeName string
 // Function that creates an upgrade handler
 CreateUpgradeHandler func(mm *module.Manager, configurator module.Configurator, keepers *keepers.AppKeepers) upgradetypes.UpgradeHandler
 // Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
 StoreUpgrades store.StoreUpgrades
}

type Fork struct {
 // Upgrade version name, for the upgrade handler, e.g. `v7`
 UpgradeName string
 // height the upgrade occurs at
 UpgradeHeight int64

 // Function that runs some custom state transition code at the beginning of a fork.
 BeginForkLogic func(ctx sdk.Context, keppers *keepers.AppKeepers, mm *module.Manager)
}
```
