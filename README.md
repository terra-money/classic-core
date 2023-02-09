## What is Terra Classic?

**[Terra](https://terra.money)** is a public, open-source blockchain protocol that provides fundamental infrastructure for a decentralized economy and enables open participation in the creation of new financial primitives to power the innovation of money.


**Classic** is the reference implementation of the Terra protocol, written in Golang. Terra Core is built atop [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and uses [Tendermint](https://github.com/tendermint/tendermint) BFT consensus. If you intend to work on Terra Core source, it is recommended that you familiarize yourself with the concepts in those projects.

Upon the implosion of Terra, a group of rebels seized control of the blockchain.  Terra's future is uncertain, but the rebels are now firmly in control. 

## Installation

### Binaries

The easiest way to get started is by downloading a pre-built binary for your operating system. You can find the latest binaries on the [releases](https://github.com/classic-terra/core/releases) page.

### From Source

**Step 1. Install Golang**

Go v1.18 is required for Terra Core.

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

**Step 2: Get Terra Core source code**

Use `git` to retrieve Terra Core from the [official repo](https://github.com/terra-money/core/) and checkout the `main` branch. This branch contains the latest stable release, which will install the `terrad` binary.

```bash
git clone https://github.com/terra-money/core/
cd core
git checkout main
```

**Step 3: Build Terra core**

Run the following command to install the executable `terrad` to your `GOPATH` and build Terra Core. `terrad` is the node daemon and CLI for interacting with a Terra node.

```bash
# COSMOS_BUILD_OPTIONS=rocksdb make install
make install
```

**Step 4: Verify your installation**

Verify that you've installed terrad successfully by running the following command:

```bash
terrad version --long
```

If terrad is installed correctly, the following information is returned:

```bash
name: terra
server_name: terrad
version: 1.0.5
commit: 8bb56e9919ecf5234a3239a6a351b509451f9d5d
build_tags: netgo,ledger
go: go version go1.18.1 linux/amd64
```

## `terrad`

**NOTE:** `terracli` has been deprecated and all of its functionalities have been merged into `terrad`.

`terrad` is the all-in-one command for operating and interacting with a running Terra node. For comprehensive coverage on each of the available functions, see [the terrad reference information](https://docs.terra.money/docs/develop/how-to/terrad/README.html). To view various subcommands and their expected arguments, use the `$ terrad --help` command:

<pre>
        <div align="left">
        <b>$ terrad --help</b>

        Stargate Terra App

        Usage:
          terrad [command]

        Available Commands:
          add-genesis-account Add a genesis account to genesis.json
          collect-gentxs      Collect genesis txs and output a genesis.json file
          debug               Tool for helping with debugging your application
          export              Export state to JSON
          gentx               Generate a genesis tx carrying a self delegation
          help                Help about any command
          init                Initialize private validator, p2p, genesis, and application configuration files
          keys                Manage your application's keys
          migrate             Migrate genesis to a specified target version
          query               Querying subcommands
          rosetta             spin up a rosetta server
          start               Run the full node
          status              Query remote node for status
          tendermint          Tendermint subcommands
          testnet             Initialize files for a terrad testnet
          tx                  Transactions subcommands
          unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
          validate-genesis    validates the genesis file at the default location or at the location passed as an arg
          version             Print the application binary version information

        Flags:
          -h, --help                help for terrad
              --home string         directory for config and data (default "/Users/$HOME/.terra")
              --log_format string   The logging format (json|plain) (default "plain")
              --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
              --trace               print out full stack trace on errors

        <b>Use "terrad [command] --help" for more information about a command.</b>
        </div>
</pre>

## Node Setup

Once you have `terrad` installed, you will need to set up your node to be part of the network.

### Join the mainnet

The following requirements are recommended for running a `columbus-5` mainnet node:

- **4 or more** CPU cores
- At least **2TB** of disk storage
- At least **100mbps** network bandwidth
- An Linux distribution

For configuration and migration instructions for setting up a Columbus-5 mainnet node, visit [The mainnet repo](https://github.com/terra-money/mainnet).

**Terra Node Quick Start**
```
terrad init nodename
wget -O ~/.terra/config/genesis.json https://cloudflare-ipfs.com/ipfs/QmZAMcdu85Qr8saFuNpL9VaxVqqLGWNAs72RVFhchL9jWs
curl https://network.terra.dev/addrbook.json > ~/.terrad/config/addrbook.json
terrad start
```
