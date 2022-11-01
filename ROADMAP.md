# Rebel Roadmap

Historical releases of this chain, before the rebel takeover, can be found at https://github.com/terra-money/classic-core


# v1.0.2: Released and running

* Declare rebellion
* Move official repository
* patch for dragonberry


## v2.0.0: Draft release out 1/11/2022, validators should audit

* upgrade to cosmos-sdk v0.45.10
* remain on patched tendermint v0.34.14 due to changes to mempool prioritization upstream
* upgrade to iavl v0.19.4
* uprage to ibc-go v1.5.0
* use a pebbledb enabled tm-db by default (goleveldb is still default here)
* **reenable ibc to cosmos chains**
* code linting, fixes by Notional
* No changes to wasm contracts
* no changes to oracle style


## v3.0.0: Address mempool prioritization issues

* upgrade to tendermint v0.34.22 and remove dependencies on patches to tendermint

## v4.0.0: branch beginning soon

* upgrade to cosmos-sdk v0.46.x or v0.47.x
* upgrade to ibc-go v5.0.1 or later
* upgrade to wasmd v0.29.1-46 by Notional (funded by Juno) or later
* ensure that Oracle uses tendermint consensus native prioritization
* Enable IBC contracts
