<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->
# Changelog
## [v2.1.2](https://github.com/classic-terra/core/releases/tag/v2.1.2) 

### Bug fixes
* (wasm) [#259](https://github.com/classic-terra/core/pull/259) add legacy wasm types
* [#268](https://github.com/classic-terra/core/pull/268) revert "feat(feeshare): Developer Incentives (#144)"

## [v2.1.1](https://github.com/classic-terra/core/releases/tag/v2.1.1) - Jun 10, 2023
### State Machine Breaking
* Wasmvm v1.1.2 (Cherry patch included)
* Wasmd [v0.30.0-terra.2](https://github.com/classic-terra/wasmd/releases/tag/v0.30.0-terra.2)
* IBC v4.3.1 (Huckleberry patch included)
* Minimum 5% commission on all validators

### Bug Fixes
* Infinite feegrant bug fixed

## [v2.0.1](https://github.com/classic-terra/core/releases/tag/v2.0.1) - May 9, 2023
### State Machine Breaking
* Apply CosmWasm security patch: Codename Cherry

## [v2.0.0](https://github.com/classic-terra/core/releases/tag/v2.0.1) - May 5, 2023

### State Machine Breaking
* (app) [#176](https://github.com/classic-terra/core/pull/176) CosmosSDK v0.45.13
* (app) [#177](https://github.com/classic-terra/core/pull/177) V3 Upgrade handler
* (gov) [#138 #185](https://github.com/classic-terra/core/pull/138) Minimum initial deposit

## Improvements
* [#172](https://github.com/classic-terra/core/pull/172) Terra operator docker refactor
* [#194](https://github.com/classic-terra/core/pull/194) Testnet refactor
* [#168](https://github.com/classic-terra/core/pull/168) refactor app folder and add 'make lint-strict'
* [#167](https://github.com/classic-terra/core/pull/167) Test framework implementation plan

## [v1.1.0](https://github.com/classic-terra/core/releases/tag/v1.1.0) - Feb 27, 2023
### Features

* (build) [#101](https://github.com/classic-terra/core/pull/101) Upgrade test
* (ante) [#103 #113 #134](https://github.com/classic-terra/core/pull/103) Add burn tax split logic
* (ante) [#107 #137 #149](https://github.com/classic-terra/core/pull/107) Burn Tax Whitelist
* (app) [#128](https://github.com/classic-terra/core/pull/128) Panic at InitChainer for the Columbus mainnet

### Improvements
* (build) [#93](https://github.com/classic-terra/core/pull/93) Use golang 1.18 and fix ad-hoc security vulnerabilities
* (build) [#97](https://github.com/classic-terra/core/pull/97) Change module path to classic-terra/core
* (build) [#102](https://github.com/classic-terra/core/pull/102) Snyk secops patches
* (build) [#105](https://github.com/classic-terra/core/pull/105) Update docker assets
* (build) [#118](https://github.com/classic-terra/core/pull/118) localnet for Apple Silicon
* [#102](https://github.com/classic-terra/core/pull/112) Upgrade SDK proto to v0.7

### Bug Fixes
* (auth/client) [#106](https://github.com/classic-terra/core/pull/106) fix ungraceful error on failed client tax query

## [v1.0.5-full-archive](https://github.com/classic-terra/core/releases/tag/v1.0.5-full-archive) - 2023-02-03

### Bug Fixes
* (app) [9c2c556](https://github.com/classic-terra/core/commit/9c2c5566dddf4459ff599d3a9ffd4ab96a2df1d0) Fix indexing bug and sync from columbus-5 genesis

## [v1.0.5](https://github.com/classic-terra/core/releases/tag/v1.0.5) - 2023-01-18

### State Machine Breaking
* (app) [#44](https://github.com/classic-terra/core/pull/44/files) Move version map register to begin block

### Bug Fixes
* (x/auth) [ebba05](https://github.com/classic-terra/core/commit/ebba0521fec4fc5655d90c0b3fdb2dbb2ec8d11f) Added Z's fix for the feeutils.go
* (x/auth) [7fe446](https://github.com/classic-terra/core/commit/7fe4468fab7a767b8779e093d671a69f26b19781) allow uluna to be taxed

## [v1.0.4](https://github.com/terra-money/classic-core/releases/tag/v0.5.23) - 2022-12-07

### State Machine Breaking
* [#800](https://github.com/terra-money/classic-core/pull/800) Tax and Burn
* [#800](https://github.com/terra-money/classic-core/pull/800) Re-enable IBC

### Bug Fixes
* Dragonberry upgrade to v0.44.5-patch and ics23 replace

## [v0.5.22](https://github.com/terra-money/classic-core/releases/tag/v0.5.22) - 2022-08-17

### State Machine Breaking
* (x/staking) [#104](https://github.com/terra-money/cosmos-sdk/pull/104) Move validator voting restrictions from Msg handler to Keeper by @edk208

## [v0.5.21](https://github.com/terra-money/classic-core/releases/tag/v0.5.21) - 2022-08-10

### State Machine Breaking
* [#784](https://github.com/terra-money/classic-core/pull/784) Burn Tax via utilization of existing Stability Tax code by @edk208

## [v0.5.20](https://github.com/terra-money/classic-core/releases/tag/v0.5.20) - 2022-05-13

### State Machine Breaking

* [#760](https://github.com/terra-money/classic-core/pull/760) Make min spread to 100% to disable swap
* [#760](https://github.com/terra-money/classic-core/pull/760) Disable IBC channels for Osmosis, Crecent and Juno

## [v0.5.19](https://github.com/terra-money/classic-core/releases/tag/v0.5.19) - 2022-05-13

### State Machine Breaking
* Disable staking power change for protecting the network

## [v0.5.18](https://github.com/terra-money/classic-core/releases/tag/v0.5.18) - 2022-04-14
This release contains a wasmer version bump from v2.0.0 to v2.2.1. The wasm caches of these two versions are not compatible, thus rebuilding is required.
To avoid possible sync delays due to the runtime rebuilding overhead, it is highly recommended that node operators rebuild their wasm cache with the [cosmwasm-cache-rebuilder](https://github.com/terra-money/cosmwasm-cache-rebuilder) before replacing terrad runtime to v0.5.18.

**Node upgrade instructions**
1. Rebuild your wasm cache using the [cosmwasm-cache-rebuilder](https://github.com/terra-money/cosmwasm-cache-rebuilder). This rebuilder can be run simultaneously without killing a running terrad process.
2. You can ignore the _file already open_ error. If other errors occur (no disk space, etc), it is safe to run the rebuilder multiple times.
3. When the rebuilder is finished, update terrad to v0.5.18 and restart.
4. The rebuilder creates a `$TERRA_HOME/data/wasm/cache/modules/v3-wasmer1` directory. The `$TERRA_HOME/data/wasm/cache/modules/v1` can be deleted after updating terrad to v0.5.18.

**Additional notes**
- :tada: Wasmer v2.2 supports the much anticipated Apple M1 processor! (See: https://wasmer.io/posts/wasmer-2.2) This means nodes can be run on any arm64 system, such as Apple Silicon machines and AWS Graviton2 Arm64.
- [\#690](https://github.com/terra-money/classic-core/pull/690) Fix wasm memory leak and module drop order. (edited) 

## [v0.5.17](https://github.com/terra-money/classic-core/releases/tag/v0.5.17) - 2022-02-22

This release contains security update.

Details will be shared after checking the status of the validator update.

## [v0.5.16](https://github.com/terra-money/classic-core/releases/tag/v0.5.16) - 2022-02-06

### State Machine Breaking
- [\#674](https://github.com/terra-money/classic-core/pull/674) fix to charge gas cost for the key length in wasm execution

#### Softfork Schedule

|               | columbus-5                | bombay-12                 |
| ------------- | ------------------------- | ------------------------- |
| Height        | 6,470,000                 | 7,800,000                 |
| Expected Time | 2022-02-13T15:00:00 (UTC) | 2022-02-11T15:00:00 (UTC) |

## [v0.5.15](https://github.com/terra-money/classic-core/releases/tag/v0.5.15) - 2022-02-03

This release contains mainly overflow checking enhancement including custom dependency updates.

- [\#670](https://github.com/terra-money/classic-core/pull/670) audit fix (overflow check)
- [\#658](https://github.com/terra-money/classic-core/pull/658) cosmos-sdk proto update
- [\#655](https://github.com/terra-money/classic-core/pull/655) use custom dependency for cosmos-sdk and tendermint

## [v0.5.14](https://github.com/terra-money/classic-core/releases/tag/v0.5.14) - 2021-12-23

This release contains wasmvm bugfix to prevent non-deterministic error messages comes out from wasmvm.

- [\#635](https://github.com/terra-money/classic-core/pull/635) bump wasmvm to v0.16.3

## [v0.5.13](https://github.com/terra-money/classic-core/releases/tag/v0.5.13) - 2021-12-17

This release contains oracle tally bugfix, which was not sorting cross rate ballot before computing weighted median. It makes the cross rate selected quite randomly...

This update will sort cross rate ballot before computing weighted median value.

- [b073fb1](https://github.com/terra-money/classic-core/commit/b073fb18e5107caee077f214b5fad4de53b055f5) sort cross rate ballot and enforce sorted before computing weighted median

## [v0.5.12](https://github.com/terra-money/classic-core/releases/tag/v0.5.12) - 2021-12-06

This release contains [IAVL bugfix](https://github.com/cosmos/iavl/pull/449)

### Improvements
- [\#620](https://github.com/terra-money/classic-core/pull/620) bump SDK to v0.44.5

### Bug Fixes
- [\#621](https://github.com/terra-money/classic-core/pull/621) fix iterator key on export for zero height
- [\#607](https://github.com/terra-money/classic-core/pull/607) skip GetBytes when exporting not migrated code

## [v0.5.11](https://github.com/terra-money/classic-core/releases/tag/v0.5.11) - 2021-11-10

### Bug Fixes
- [\#603](https://github.com/terra-money/classic-core/pull/603) bump CosmWasm/wasmvm to v0.16.2 to fix cache position error

## [v0.5.10](https://github.com/terra-money/classic-core/releases/tag/v0.5.10) - 2021-11-03

### Improvements
- [\#601](https://github.com/terra-money/classic-core/pull/601) revert readvm-pool feature
- [\#594](https://github.com/terra-money/classic-core/pull/594) bump SDK to v0.44.3 and Tendermint to v0.34.14
- [\#593](https://github.com/terra-money/classic-core/pull/593) revert jemalloc integration from wasmvm & add troubleshoot section to docs

## [v0.5.9](https://github.com/terra-money/classic-core/releases/tag/v0.5.9) - 2021-10-22

### Bug Fixes
- [\#588](https://github.com/terra-money/classic-core/pull/588) - fix wasm external querier to use single wasmvm with context value

## [v0.5.8](https://github.com/terra-money/classic-core/releases/tag/v0.5.8) - 2021-10-21

### Bug Fixes
- [\#584](https://github.com/terra-money/classic-core/pull/584) - security update
- [\#586](https://github.com/terra-money/classic-core/pull/586) - memory leak fix
- [\#583](https://github.com/terra-money/classic-core/pull/583) - ibc ante handler to prevent duplicated packet relaying

## [v0.5.7](https://github.com/terra-money/classic-core/releases/tag/v0.5.7) - 2021-10-13

This release mainly contains CosmosSDK security updates

### Bug Fixes
- [\#578](https://github.com/terra-money/classic-core/pull/578) Bump CosmosSDK to v0.44.2
- [\#576](https://github.com/terra-money/classic-core/pull/576) Change write-vm-memory-cache-size default config replace location
- [\#573](https://github.com/terra-money/classic-core/pull/573) Fix to distribute oracle rewards for all whitelist denoms

## [v0.5.6](https://github.com/terra-money/classic-core/releases/tag/v0.5.6) - 2021-10-06

This release contains updates for multi-reader thread implementation with necessary dependency updates for multi-reader thread implementation.

* Bump cosmos-sdk to [v0.44.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.1)
* Bump tendermint to [v0.44.1](https://github.com/tendermint/tendermint/releases/tag/v0.34.13)
* Bump CosmWasm to [v0.16.2](https://github.com/CosmWasm/cosmwasm/releases/tag/v0.16.2)

### Improvements
- [\#546](https://github.com/terra-money/classic-core/pull/546) Implement read VM pool

## [v0.5.5](https://github.com/terra-money/classic-core/releases/tag/v0.5.5) - 2021-09-27

### Improvements
- [\#564](https://github.com/terra-money/classic-core/pull/564) Legacy tx encode support

### Bug Fixes
- [\#566](https://github.com/terra-money/classic-core/pull/566) Emit missing wasm tax events.

## [v0.5.4](https://github.com/terra-money/classic-core/releases/tag/v0.5.4) - 2021-09-24

### Bug Fixes
- [\#560](https://github.com/terra-money/classic-core/pull/560) Fix migration bug of multisig pubkey which was in v040 auth module migration.

## [v0.5.3](https://github.com/terra-money/classic-core/releases/tag/v0.5.3) - 2021-09-20

### Improvement
- [\#551](https://github.com/terra-money/classic-core/pull/551) Bump CosmosSDK to [v0.44.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.0) and IBC-GO to [v1.1.0](https://github.com/cosmos/ibc-go/releases/tag/v1.1.0) both for security patch

### Bug Fixes
- [\#544](https://github.com/terra-money/classic-core/pull/544) update rest&cli interface to receive base64 key for wasm raw querier
- [\#557](https://github.com/terra-money/classic-core/pull/557) fix proto urls

## [v0.5.2](https://github.com/terra-money/classic-core/releases/tag/v0.5.2) - 2021-08-26
- [\#539](https://github.com/terra-money/classic-core/pull/539) Bump ibc-go version to v1.0.1 to fix security vulnerability

## [v0.5.1](https://github.com/terra-money/classic-core/releases/tag/v0.5.1) - 2021-08-18

### Bug Fixes
- [\#536](https://github.com/terra-money/classic-core/pull/536) Sort migrated map array with key to generate deterministic result

## [v0.5.0](https://github.com/terra-money/classic-core/releases/tag/v0.5.0) - 2021-08-11

### Improvements

- [\#529](https://github.com/terra-money/classic-core/pull/529) Bump CosmosSDK to [v0.43.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.43.0)
- [\#523](https://github.com/terra-money/classic-core/pull/526) Bump CosmosSDK to [v0.43.0-rc3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.43.0-rc3) and CosmWasm to [v0.16.0](https://github.com/CosmWasm/cosmwasm/release/tag/v0.16.0)
- [\#520](https://github.com/terra-money/classic-core/pull/520) Implement ContractInfo wasm query interface for the contracts
- [\#518](https://github.com/terra-money/classic-core/pull/518) Rollback market base pool split
- [\#507](https://github.com/terra-money/classic-core/pull/507) Bump SDK to [v0.43.0-rc0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.43.0-rc0).
- [\#507](https://github.com/terra-money/classic-core/pull/507) Bump Tendermint to [v0.34.11](https://github.com/tendermint/tendermint/releases/tag/v0.34.11).
- [\#509](https://github.com/terra-money/classic-core/pull/509) Bump cosmwasm to [v0.15.1](https://github.com/CosmWasm/cosmwasm/release/tag/v0.15.1).
- [\#467](https://github.com/terra-money/classic-core/pull/467) Split market base pool into mint & burn base pool
- [\#465](https://github.com/terra-money/classic-core/pull/465) [\#472](https://github.com/terra-money/classic-core/pull/472)  Bump cosmwasm to [v0.14.0](https://github.com/CosmWasm/cosmwasm/release/tag/v0.14.0).
- [\#477](https://github.com/terra-money/classic-core/pull/477) Burn `reward_weight` ratio of seigniorage instead of redirecting to oracle rewards
- [260b41a](https://github.com/terra-money/classic-core/blob/471cd878b11b21f10511e272621b850a690eb6f3/x/market/keeper/msg_server.go#L121-L127) Redirect swap fees to oracle rewards (distributed over 3 years)
- Bump golang prerequisite from 1.15 to 1.16.

### Bug Fixes
- [\#497](https://github.com/terra-money/classic-core/pull/497) Fix wasm oracle querier to return Base/Quote exchange rate
- [\#529](https://github.com/terra-money/classic-core/pull/529) Fix keyring kwallet backend when using with empty wallet.

## v0.4.5
This release is a hotfix for two high-severity issues.

[Upgrade Instructions](https://github.com/terra-money/classic-mainnet/wiki/Columbus-4-Softfork-Instructions)

```
$ git fetch --all --tags
$ git checkout v0.4.5
$ make install
```

### Upgrade Time

#### Mainnet
**Target Height**: 2,380,000

```
Tue Mar 30 2021 09:00:00 GMT+0000 (UTC)
Tue Mar 30 2021 01:00:00 GMT-0800 (PST)
Tue Mar 30 2021 18:00:00 GMT+0900 (KST)
```

#### Testnet
**Target Height**: 3,150,000

```
Tue Mar 25 2021 09:00:00 GMT+0000 (UTC)
Tue Mar 25 2021 01:00:00 GMT-0800 (PST)
Tue Mar 25 2021 18:00:00 GMT+0900 (KST)
```

### Improvements
* [\#463](https://github.com/terra-money/classic-core/pull/463) Lower oracle feeder cost.
* [\#462](https://github.com/terra-money/classic-core/pull/462) Softfork to cap the max tx limit with ConsensusParam update.

## v0.4.4

### Bug Fixes
* [\#461](https://github.com/terra-money/classic-core/pull/461) Update Dockerfile script to use proper libgo_cosmwasm_musla.a file

## v0.4.3

### Improvements
* [\#460](https://github.com/terra-money/classic-core/pull/460) Add `tx-gas-hard-limit` flag to filter out tx with abnormally huge gas

### Bug Fixes
* [\#456](https://github.com/terra-money/classic-core/pull/456) `RewardWeightUpdateProposal` CLI parse error
* [\#454](https://github.com/terra-money/classic-core/pull/454) `MsgSwapSend` rest interface parse error

## v0.4.2

### Release Note

This release is a hotfix for two high-severity issues in the currently live Terra Core@0.4.1.
* [\#440](https://github.com/terra-money/classic-core/pull/440) go-cosmwasm iterator memory leak
* [\#445](https://github.com/terra-money/classic-core/pull/445) treasury division by zero protection


### How to Upgrade
You can stop, update and restart terrad anytime before the upgrade time.

[Upgrade Instructions](https://github.com/terra-money/classic-mainnet/wiki/Columbus-4-Hotfix-(v0.4.2)-Instructions)

```
$ git fetch --all --tags
$ git checkout v0.4.2
$ make install
```

### Upgrade Time
**Target Height**: 1,915,199

```
Tue Feb 23 2021 02:27:50 GMT+0000 (UTC)
Tue Feb 23 2021 11:27:50 GMT+0900 (KST)
Mon Feb 22 2021 18:27:50 GMT-0800 (PST)
```

## v0.4.1

### Release Notes
**This upgrade contains softfork**
Please understand the details and apply it before the target height.

[Upgrade Details](https://classic-agora.terra.money/t/terra-core-v0-4-1-soft-fork-upgrade-recommendation/262)
[Upgrade Instructions](https://github.com/terra-money/mainnet/wiki/Columbus-4-Softfork-Instructions)

### How to Upgrade
It is softfork, so you can update terrad anytime before the upgrade time.
```
$ git fetch --all --tags
$ git checkout v0.4.1
$ make install
```

### Upgrade Time

* **Target Height for `columbus-4`**: 1200000
* **Target Height for `tequila-0004`**: 1350000

```
// MAINNET
// Fri Jan 01 2021 18:00:00 GMT+0900 (KST)
// Fri Jan 01 2021 09:00:00 GMT+0000 (UTC)
// Fri Jan 01 2021 01:00:00 GMT-0800 (PST)
//
// TEQUILA
// Fri Nov 27 2020 12:00:00 GMT+0900 (KST)
// Fri Nov 27 2020 03:00:00 GMT+0000 (UTC)
// Thu Nov 26 2020 19:00:00 GMT-0800 (PST)
```

### Improvements
* [\#426](https://github.com/terra-money/classic-core/pull/426) CosmWasm Cache Implementation (100x faster than before)
* [\#413](https://github.com/terra-money/classic-core/pull/413) CosmWasm Logging Whitelist

### Bug Fixes
* [\#427](https://github.com/terra-money/classic-core/pull/427) CosmWasm Staking Query

### Param Changes
* [\#433](https://github.com/terra-money/classic-core/pull/433) Increase ExecuteMsgSize limit to 4096 from 1024

## v0.4.0

### Release Notes
- [Cosmos-SDK v0.38 Release Notes](https://github.com/cosmos/cosmos-sdk/wiki/v0.38-Release-Notes)
- [Cosmos-SDK v0.39.0 Release Notes](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.39.0)
- [Cosmos-SDK v0.39.1 Release Notes](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.39.1)
- [Cosmos-SDK Breaking Changes](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.38.0)

### Improvements

* [\#407](https://github.com/terra-money/classic-core/pull/407) Allow `gov/MsgVote` to be grantable
* [\#405](https://github.com/terra-money/classic-core/pull/405) CosmWasm oracle exchange rates query interface
* [\#388](https://github.com/terra-money/classic-core/pull/388) Bump CosmWasm to v0.10.1
* [\#383](https://github.com/terra-money/classic-core/pull/383) Bump SDK version to v0.39.1
* [\#374](https://github.com/terra-money/classic-core/pull/374) Bump SDK version to v0.39 and CosmWasm to v0.9.4
* [\#357](https://github.com/terra-money/classic-core/pull/357) Bump CosmWasm to v0.9
* [\#352](https://github.com/terra-money/classic-core/pull/352) MsgAuthorization module to allow subkey feature
* [\#349](https://github.com/terra-money/classic-core/pull/349) Add `--old-hd-path` flag to support 118 coin type users
* [\#348](https://github.com/terra-money/classic-core/pull/348) MsgSwapSend to allow sending all swap coin
* [\#347](https://github.com/terra-money/classic-core/pull/347) CosmWasm custom msg & querier handler
* [\#343](https://github.com/terra-money/classic-core/pull/343) Burn Address
* [\#335](https://github.com/terra-money/classic-core/pull/335) CosmWasm integration
* [\#325](https://github.com/terra-money/classic-core/pull/325) New oracle msgs for vote process optimization
* [\#324](https://github.com/terra-money/classic-core/pull/324) Update to emit events at proposal handler
* [\#323](https://github.com/terra-money/classic-core/pull/323) Bump SDK version to v0.38.x

### Bug Fixes
* [\#360](https://github.com/terra-money/classic-core/pull/360) Fix market module pool adjustment to apply delta with actual minted amount
* [\#336](https://github.com/terra-money/classic-core/pull/336) Allow zero tobin tax rate

### Breaking Changes

#### Keys Migration
Any existing keys that were managed via Keybase in prior versions must be migrated. To migrate keys, execute the following:
```
$ terracli keys migrate [--home] [--keyring-backend]
```

The above command will provide a prompt for each existing key and ask if you wish for it to be skipped or not. If the key is not to be skipped, you must provide the correct passphrase for it to be migrated successfully.

#### Pruning Configuration

The operator can now set the pruning options by passing a pruning configuration via command line option or `app.toml`. The pruning flag supports the following
options: `default`, `everything`, `nothing`, `custom` - see the [PR](https://github.com/cosmos/cosmos-sdk/pull/6475) for further details. If the operator chooses `custom`, they may want to provide either of the granular pruning values:

- `pruning-keep-recent`
- `pruning-keep-every`
- `pruning-interval`

The former two options dictate how many recent versions are kept on disk and the offset of what versions are kept after that
respectively, and the latter defines the height interval in which versions are deleted in a batch.

**The operator, who wants to upgrade the node from v0.3 to v0.4, must change pruning option in `app.toml` to one of above options.**

#### API Changes
* The `block_meta` field has been removed from `/blocks/{block_height}` becasuse it was redandunt data with `block_header`.
* The `whitelist` of`/oracle/parameters` response has been changed from `[]string` to `[]{ name: string; tobin_tax: string; }`

## v0.3.6

### Improvements
#### [99581ba](https://github.com/terra-money/classic-core/commit/99581baf89a838cf09a25d47adc2fd2cc97ab4a2) Ledger update(custom ledger library) & Bump SDK to v0.37.13

## v0.3.5

### Improvements
#### [654b5cb](https://github.com/terra-money/classic-core/commit/654b5cb66a9152dcf6e53f73e7935522251a1ede) Bump SDK to v0.37.11

### Bug Fixes
#### [7a3d01c](https://github.com/terra-money/classic-core/commit/7a3d01c9198cfdcc67d90593c92ce5cb465e4516) Oracle slashing unbonding state check

## v0.3.4

### Improvements
#### [\#338](https://github.com/terra-money/classic-core/pull/338) Bump SDK to v0.37.9 for Tendermint security patch

## v0.3.3

### Improvements
#### [\#319](https://github.com/terra-money/classic-core/pull/319) Bump SDK to v0.37.6
#### [\#321](https://github.com/terra-money/classic-core/pull/321) Revert to distribute zero oracle rewards

## v0.3.2

### Improvements
#### [\#313](https://github.com/terra-money/classic-core/pull/313) upgrade SDK
* Bump SDK version to [v0.37.5](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.37.5)
* Tendermint version to [v0.32.8](https://github.com/tendermint/tendermint/releases/tag/v0.32.8)
#### [\#312](https://github.com/terra-money/classic-core/pull/312) upgrade golangci-lint version to v1.22.2

## v0.3.1

### Bug Fixes
#### [\#303](https://github.com/terra-money/classic-core/pull/303) fix estimate fee endpoint for multiple signature tx
#### [\#304](https://github.com/terra-money/classic-core/pull/304) genesis scrpit update

### Improvements
#### [\#301](https://github.com/terra-money/classic-core/pull/301) README update
#### [\#305](https://github.com/terra-money/classic-core/pull/305) swagger update
#### [\#306](https://github.com/terra-money/classic-core/pull/306) circleci update for goreleaser

## v0.3.0
### Breaking Changes
#### [\#265](https://github.com/terra-money/classic-core/pull/265) Oracle refactor & Oracle slashing
##### Slashing
A validator get slashed `SlashFraction`% if the one perform any of the following violations in `SlashWindow - minValidPerWindow` voteperiods over a window of `SlashWindow` voteperiods:

1. A vote is missing for any of the denom in the whitelist. Oracle voters looking to abstain must still submit a "vote of no confidence", which has 0 for the luna ExchangeRate field of the prevote.
2. A submitted vote is more than max(`RewardBand`, standard deviation) from the elected median
3. Oracle voters who submit abstain vote with other invalid votes will also get slashed.

##### Codec Changes
* `oracle/MsgDelegateFeederPermission` => `oracle/MsgDelegateFeedConsent`

##### New End Points
* `/oracle/voters/{%s}/miss` return the # of vote periods missed in this oracle slash window.

##### Path Changes
* `/oracle/denoms/{denom}/price` => `/oracle/denoms/{denom}/exchange_rate`

##### Request Body Changes

* POST /oracle/denoms/{denom}/prevotes
* POST /oracle/denoms/{denom}/votes
```
Price sdk.Dec `json:"price"`
```
has been changed to
```
ExchangeRate sdk.Dec `json:"exchange_rate"`
```

#### [\#256](https://github.com/terra-money/classic-core/pull/256) Oracle endpoints improvement
##### New EndPoints
```
/oracle/voters/{validator}/votes
/oracle/voters/{validator}/prevotes
/oracle/denoms/prices
```

#### [\#250](https://github.com/terra-money/classic-core/pull/250) Oracle whitelist & Reward distribution update
* Create a whitelist param that stores an array of denoms that are whitelisted by the protocol.
* Edit the oracle `Reward Pool of a VotePeriod = oracle module account / (n vote periods)`.
* Oracle module account is whitelisted in the bank module such that users can donate funds to the oracle module account

#### [\#234](https://github.com/terra-money/classic-core/pull/234) Adopt gov module
`distribution` module already contains `community-pool-spend` proposal suitable for `budget` so budget module is removed. There are two custom governance proposals from `treasury` module; `tax-rate-update` & `reward-weight-update` proposals.

##### New EndPoints
```
(GET)/gov/proposals
(GET)/gov/proposals/{proposalId}
(GET)/gov/proposals/{proposalId}/proposer
(GET)/gov/proposals/{proposalId}/deposits
(GET)/gov/proposals/{proposalId}/deposits/{depositor}
(GET/POST)/gov/proposals/{proposalId}/votes
(GET)/gov/proposals/{proposalId}/votes/{voter}
(GET)/gov/proposals/{proposalId}/tally
(GET)/gov/parameters/deposit
(GET)/gov/parameters/tallying
(GET)/gov/parameters/voting
(POST)/gov/proposals/tax_rate_update
(POST)/gov/proposals/reward_weight_update
(POST)/gov/proposals/param_change
(POST)/gov/proposals/community_pool_spend
```

#### [\#233](https://github.com/terra-money/classic-core/pull/233) Swap constant product
As proposed [here](https://classic-agora.terra.money/uploads/short-url/92QHxFtEmWUEwf9kWTminuobwpM.pdf), apply constant product to swap feature.

##### Compute Pools
```
// Both LUNA and TERRA pools are using SDR units.
cp = basePool*basePool
terraPool = (basePool + terraDelta)
lunaPool = cp/terraPool
```

##### LUNA to TERRA swap
```
// offerAmt must be SDR units
newLunaPool = lunaPool + offerAmt
newTerraPool = cp / newLunaPool
returnAmt = newTerraPool - terraPool

// Swap return SDR Amt to TERRA
returnLunaAmt = market.swap(returnAmt, "LUNA")
```

##### TERRA to LUNA swap
```
// offerAmt must be SDR units
newTerraPool = terraPool + offerAmt
newLunaPool = cp / newTerraPool
returnAmt = newLunaPool - lunaPool

// Swap return SDR Amt to proper TERRA
returnTerraAmt = market.swap(returnAmt, "TERRA")
```

##### TERRA to TERRA swap

Apply only fixed tobin-tax without computing and changing pools


##### New EndPoints
```
/market/terra_pool_delta
```

#### [\#231](https://github.com/terra-money/classic-core/pull/231) Bump SDK to v0.37.x
##### REST end points, which are changed
All REST responses now wrap the original resource/result. The response
will contain two fields: height and result.
```
/market/params => /market/parameters
/oracle/params => /oracle/parameters
/treasury/tax-rate => /treasury/tax_rate
/treasury/tax-rate/{epoch} => /treasury/tax_rate/{epoch}
/treasury/tax-cap => /treasury/tax_cap
/treasury/tax-cap/{denom} => /treasury/tax_cap/{denom}
/treasury/reward-weight => /treasury/reward_weight
/treasury/reward-weight/{epoch} => /treasury/reward_weight/{epoch}
/treasury/tax-proceeds => /treasury/tax_proceeds
/treasury/tax-proceeds/{epoch} => /treasury/tax_proceeds/{epoch}
/treasury/seigniorage-proceeds => /treasury/seigniorage_proceeds
/treasury/seigniorage-proceeds/{epoch} => /treasury/seigniorage_proceeds/{epoch}
/treasury/current-epoch => /treasury/current_epoch
/treasury/params => /treasury/parameters
```
##### REST end points, which response object key is removed
```
/treasury/current_epoch
/treasury/seigniorage_proceeds/{epoch}
/treasury/seigniorage_proceeds
/treasury/tax_proceeds/{epoch}
/treasury/tax_proceeds
/treasury/reward_weight/{epoch}
/treasury/reward_weight
/treasury/tax_cap/{denom}
/treasury/tax_rate/{epoch}
/treasury/tax_rate
/oracle/denoms/actives
/oracle/denoms/{denom}/price
/oracle/denoms/{denom}/prevotes/{voter}
/oracle/denoms/{denom}/prevotes
/oracle/denoms/{denom}/votes/{voter}
/oracle/denoms/{denom}/votes
```
##### New REST endpoints
```
/supply/total
/supply/total/{denomination}
/market/last_day_issuance
/oracle/voters/{%s}/voting_info
/oracle/voting_infos
/treasury/historical_issuance/{epoch}
```
##### Codec changes
```
auth/Account => core/Account
auth/StdTx => core/StdTx
pay/MsgSend => bank/MsgSend
pay/MsgMultiSend => bank/MsgMultiSend
```
##### Other Changes
* GradedVestingAccount is fully removed
* LazyGradedVestingAccount's vesting_lazy_schedules is changed to vesting_schedules
* Improve the UX of fee and tax. The sender have to specify fees containing tax amount with following methods.
  1. Use `/bank/accounts/{address}/transfers` without fees. `terracli` will compute tax amount and fill fees field containing both gas & tax fee.
  2. Use `/txs/estimate_fee` to estimate fees of StdTX, and replace StdTx.Fee.Amount to estimated fee
  3. Compute tax with `/treasury/tax_rate` & `/treasury/tax_cap` add computed tax with original gas fee

## v0.2.4
### Bug fixes
#### [\#196](https://github.com/terra-money/classic-core/pull/196) peek epoch seigniorage
Change PeekEpochSeigniorage to compute seigniorage by subtracting current issuance from previous issuance

#### [\#198](https://github.com/terra-money/classic-core/pull/198) Use next block for treasury tax and reward update
updateTaxPolicy and updateRewardPolicy are updating new tax-rate and reward-weight with current ctx. The ctx height is the last block of current epoch, but treasury should update next epoch's tax-rate and reward-weight at the last block of current epoch.
In updateTaxPolicy and updateRewardPolicy, change ctx input of keeper setter to ctx with next epoch height.

### Features
#### [\#193](https://github.com/terra-money/classic-core/pull/193) Recover old hd path
Added `--old-hd-path` option to `$terracli keys add` command for recovering old bip44 path(for atom)
##### Example
```
$ terracli keys add tmp --recover --old-hd-path
Enter a passphrase to encrypt your key to disk:
Repeat the passphrase:
> Enter your bip39 mnemonic
candy hint hamster cute inquiry bright industry decide assist wedding carpet fiber arm menu machine lottery type alert fan march argue adapt recycle stomach

NAME:	TYPE:	ADDRESS:					PUBKEY:
tmp   local	terra1gaczd45crhwfa4x05k9747cuxwfmnduvmtyefs	terrapub1addwnpepqv6tse2pyag9ts5vy6dk4h3qh7xc9qhat4jx449n6nrfve3jhzldz3f3l7p
```

## v0.2.3
- [\#187](https://github.com/terra-money/classic-core/pull/187): Change all time instance timezone to UTC to remove gap in time calculation

### Changes
#### [\#187](https://github.com/terra-money/classic-core/pull/187) Bugfix/fix-time-zone
In update_230000.go, we change genesis time derivation from
```
genesisTime := time.Unix(genesisUnixTime, 0)
```
to
```
genesisTime := time.Unix(genesisUnixTime, 0).UTC()
```

## v0.2.2

- [\#185](https://github.com/terra-money/classic-core/pull/185): Improve oracle specs
- [\#184](https://github.com/terra-money/classic-core/pull/184): Fix `terracli` docs
- [\#183](https://github.com/terra-money/classic-core/pull/183): Change all GradedVestingAccounts to LazyGradedVestingAccounts.
- [\#179](https://github.com/terra-money/classic-core/pull/179): Conform querier responses to be returned in JSON format
- [\#178](https://github.com/terra-money/classic-core/pull/178): Change BIP44 PATH to 330

### Changes
#### [\#185](https://github.com/terra-money/classic-core/pull/185) Oracle `MsgFeederDelegatePermission` specs
Added docs for using `MsgFeederDelegatePermission` to oracle specs

#### [\#185](https://github.com/terra-money/classic-core/pull/185) Oracle price vote denom error fix
Oracle specs now specify micro units `uluna` and `uusd` for correct denominations for price prevotes and votes

#### [\#184](https://github.com/terra-money/classic-core/pull/184) Minor terracli fix

#### [\#183](https://github.com/terra-money/classic-core/pull/183) Oracle param update
```
OracleRewardBand: 1% => 2%
```

#### [\#183](https://github.com/terra-money/classic-core/pull/183) Market param update
```
DailyLunaDeltaCap: 0.5% => 0.1%
```

#### [\#183](https://github.com/terra-money/classic-core/pull/183) LazyGradedVestingAccount

* Spread out the cliffs for presale investors, with varying degrees of severity (details [\#180](https://github.com/terra-money/classic-core/issues/180))

#### [\#179](https://github.com/terra-money/classic-core/pull/179) Align Querier responses to JSON

* Querier was returning misaligned formats for return values, now aligned to JSON format

#### [\#178](https://github.com/terra-money/classic-core/pull/178) Correctly use 330 as the coin type field in BIP 44 PATH

* We were previously using the Cosmos coin type field for the BIP44 path. Changed to Terra's own 330.


## v0.2.1

- [\#166](https://github.com/terra-money/classic-core/pull/166): Newly added parameters were not being added to the columbus-2 genesis.json file. Fixed.

## v0.2.0

### Bug Fixes

* [\#140](https://github.com/terra-money/classic-core/pull/140) Fix export bug.

* [\#140](https://github.com/terra-money/classic-core/pull/140) Client querier bug fix (distr outstanding rewards)

* [\#140](https://github.com/terra-money/classic-core/pull/140) Fix budget module to delete all votes when submitter withdraws the program and to use DeleteVotesForProgram to delete all votes for a program.

### Improvements
#### [\#140](https://github.com/terra-money/classic-core/pull/140) Msg Types

```
cosmos-sdk/MsgSend => pay/MsgSend
cosmos-sdk/MsgMultiSend => pay/MsgMultiSend

cosmos-sdk/MsgCreateValidator => staking/MsgCreateValidator
cosmos-sdk/MsgEditValidator => staking/MsgEditValidator
cosmos-sdk/MsgDelegate => staking/MsgDelegate
cosmos-sdk/MsgUndelegate => staking/MsgUndelegate
cosmos-sdk/MsgBeginRedelegate => staking/MsgBeginRedelegate

cosmos-sdk/MsgWithdrawDelegationReward => distribution/MsgWithdrawDelegationReward
cosmos-sdk/MsgWithdrawValidatorCommission => distribution/MsgWithdrawValidatorCommission
cosmos-sdk/MsgModifyWithdrawAddress => distribution/MsgModifyWithdrawAddress

cosmos-sdk/MsgUnjail => slashing/MsgUnjail
```

#### [\#140](https://github.com/terra-money/classic-core/pull/140) Oracle updates prevoting/voting
MsgPriceFeed is split into ```MsgPricePrevote``` and ```MsgPriceVote```
```
Period  |  P1 |  P2 |  P3 |  ...    |
Prevote |  O  |  O  |  O  |  ...    |
        |-----\-----\-----\-----    |
Vote    |     |  O  |  O  |  ...    |
```
In prevote stage, a validator should submit the hash of the part of real vote msg to prove the validator is not just copying other validators price vote. In vote phrase, the validator should reveal the real price by submitting MsgPriceVote with ```salt```.

The submission order has to be kept in (vote -> prevote) order. If an prevote comes early, it will replace previous prevote so next vote, which reveals the proof for previous prevote, will be failed.

#### [\#148](https://github.com/terra-money/classic-core/pull/148) Oracle voting right delegation
By using the oracle/MsgDelegateFeederPermission a validator can assign the right to vote to another account at any time. The validator account will preserve its right to vote at any time.

#### [\#140](https://github.com/terra-money/classic-core/pull/140) & [\#148](https://github.com/terra-money/classic-core/pull/148) Rest Interface Update
##### Change rest interface url
```
"/distribution/parameters" => "/distribution/params"
"/staking/parameters" => "/staking/params"
```

##### Change request body
```
Send request body
From:
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount   sdk.Coins    `json:"coins"`
}

To:
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Coins   sdk.Coins    `json:"coins"`
}
```

##### New rest interfaces
```
(GET/POST) "/oracle/denoms/{%s}/votes"
(GET/POST) "/oracle/denoms/{%s}/prevotes"
(GET/POST) "/oracle/voters/{%s}/feeder"
(GET/POST) "/market/swap"
(GET/POST) "/market/params"
```

#### [\#140](https://github.com/terra-money/classic-core/pull/140) Add transaction logs for tax and swap amount
##### Send Tx
Add **tax** log to send transaction for recording real amount which a transaction pay.
Ex) txs/B515331BF9EA9A92AD59A85D593E5A2B170E3D297C59E85DDA9FA6FF33790E9B
```
{
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "{\"tax\":\"400uluna\"}"
    }
  ]
}
```

##### Swap Tx
Add **swap_coin** log to swap transaction for recording the amount of swapped coin along with offered coin
```
{
  "logs": [
    {
      "msg_index": 0,
      "success": true,
      "log": "{\"swap_coin\":\"400ukrw\"}"
    }
  ]
}
```

#### [\#150](https://github.com/terra-money/classic-core/pull/150) Market Swap protections

##### Add bidirectional Luna supply change cap on market swaps.
A daily trading cap (luna supply change cap) protects excessive luna volatility. Capping Luna deflation prevents divesting attacks (attacker swaps large amount into terra to avoid slippage) and consensus attacks by limiting access to staking tokens. Early parameters are 2% - 10% on both sides of the trade.

##### Add bidirectional Luna spread fees on market swaps
To protect against short term price deviations between the open market and the on-chain oracle, we now charge a 2-10% spread on swaps that involve luna.

##### Change oracle reward scheme from monthly seigniorage to validators to minute distribution
Swap spreads are distributed to oracle ballot winners on the oracle VotePeriod; this vastly shortens distribution periods. Also, all stakeholders receive oracle rewards (includes delegators).

##### Swaps halt immediately after an illiquid oracle vote
Previously we facilitated swaps for 10 VotePeriods after the last valid oracle ballot. We now stop swaps immediately to prevent arbitrage attacks from price drift.


### Parameter Changes

#### [\#150](https://github.com/terra-money/classic-core/pull/150) Change MiningRewardWeight.Max from 20% to 90%. This is to reduce volatility in fees at network infancy.

#### Changed BlocksPerMinute from 12 to 5 to more accurately reflect Columbus block times.
