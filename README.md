<p>&nbsp;</p>
<p align="center">

<img src="core_logo.svg" width=500>

</p>

<p align="center">
Full-node software implementing the Terra protocol<br/><br/>

<a href="https://codecov.io/gh/terra-money/core">
    <img src="https://codecov.io/gh/terra-money/core/branch/develop/graph/badge.svg">
</a>
<a href="https://goreportcard.com/report/github.com/terra-money/core">
    <img src="https://goreportcard.com/badge/github.com/terra-money/core">
</a>

</p>

<p align="center">
  <a href="https://docs.terra.money/"><strong>Explore the Docs »</strong></a>
  <br />
  <br/>
  <a href="https://docs.terra.money/dev">Dev Guide</a>
  ·
  <a href="https://pkg.go.dev/github.com/terra-money/core?tab=subdirectories">Go API</a>
  ·
  <a href="https://swagger.terra.money/">REST API</a>
  ·
  <a href="https://docs.terra.money/#sdks-for-developers">SDKs</a>
  ·
  <a href="https://finder.terra.money/">Finder</a>
  ·
  <a href="https://station.terra.money/">Station</a>
</p>

<br/>

## What is Terra?

**[Terra](https://terra.money)** is a blockchain protocol that provides fundamental infrastructure for a decentralized economy and enables open participation in the creation of new financial primitives to power the innovation of money.

The Terra blockchain is secured through distributed consensus over native staked asset Luna, and supports the issuance of price-tracking stablecoins (TerraKRW, TerraUSD, etc.) that are pegged to major world currencies. Smart contracts on Terra run on WebAssembly and can take advantage of core modules like on-chain swaps, price oracle, and staking rewards to power modern DeFi apps. Through fiscal policy managed by community governance, Terra is a democratized economy regulated by its users.

**Terra Core** is the reference implementation of the Terra protocol, written in Golang. Terra Core is built atop [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and uses [Tendermint](https://github.com/tendermint/tendermint) BFT consensus. If you intend to work on Terra Core source, it is recommended that you familiarize yourself with the concepts in those projects.

## Installation

### Binaries

You can find the latest binaries on our [releases](https://github.com/terra-money/core/releases) page.

### From Source

We recommend the following for running Terra Core:

- **2 or more** CPU cores
- At least **300GB** of disk storage
- At least **2.5 - 5mbps** network bandwidth

#### Step 1. Install Golang

Go v1.14+ or higher is required for Terra Core.

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

#### Step 2: Get Terra Core source code

Use `git` to retrieve Terra Core from the [official repo](https://github.com/terra-money/core/), and checkout the `master` branch, which contains the latest stable release. That should install the `terrad` binary.

```bash
git clone https://github.com/terra-money/core/
cd core
git checkout master
```

#### Step 3: Build from source

You can now build Terra Core. Running the following command will install executable `terrad` (Terra node daemon and CLI for interacting with the node) to your `GOPATH`.

```bash
make install
```

#### Step 4: Verify your installation

Verify that everything is OK. If you get something like the following, you've successfully installed Terra Core on your system.

```bash
terrad version --long
name: terra
server_name: terrad
version: 0.5.0-rc0-9-g640fd0ed
commit: 640fd0ed921d029f4d1c3d88435bd5dbd67d14cd
build_tags: netgo,ledger
go: go version go1.16.5 darwin/amd64
```

### CLI

`terrad` provides you a command-line interface to a running node, communicating over RPC. You can find comprehensive coverage on how to use the CMD on our [official docs](https://docs.terra.money/terracli). The various subcommands and their expected arguments can also be discovered by issuing:

<pre>
        <div align="left">
        <b>$ terrad --help</b>

        Command line interface for interacting with terrad

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
              --home string         directory for config and data (default "/Users/yeoyunseog/.terra")
              --log_format string   The logging format (json|plain) (default "plain")
              --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
              --trace               print out full stack trace on errors

        <b>Use "terracli [command] --help" for more information about a command.</b>
        </div>
</pre>

## Node Setup

#### Active Networks

| Chain ID       | Description        | Public Node (LCD)             |
| -------------- | ------------------ | ----------------------------- |
| `columbus-4`   | Mainnet            | https://lcd.terra.dev         |
| `tequila-0004` | Columbus-4 Testnet | https://tequila-lcd.terra.dev |
| `bombay-0008`  | Columbus-5 Testnet | https://bombay-lcd.terra.dev  |

### Running a Local Testnet

The simplest Terra network you can set up will be a local testnet with just a single node. You will create one account and be the sole validator signing blocks for the network.

#### Step 1. Create network and account

First, initialize your genesis file that will bootstrap the network. Set a name for your local testnet, and provide a moniker to refer to your node.

```bash
terrad init --chain-id=<testnet_name> <node_moniker>
```

You will need a Terra account to start. You can generate one with:

```bash
terrad keys add <account_name>
```

#### Step 2. Add account to genesis

Next, you need to add your account to the genesis. The following commands add your account and set the initial balance:

```bash
terrad add-genesis-account $(terrad keys show <account_name> -a) 100000000uluna,1000usd
terrad gentx --name my_account --amount 10000000uluna
terrad collect-gentxs
```

#### Step 3. Run Terra daemon

Now, you can start your private Terra network:

```bash
terrad start
```

Your `terrad` node should now be running a node on `tcp://localhost:26656`, listening for incoming transactions and signing blocks. You've successfully set up your local Terra network!

### Joining the mainnet

[The mainnet repo](https://github.com/terra-money/mainnet) contains snapshot of the launch as well as network updates.

### Joining a testnet

[Our testnet repo](https://github.com/terra-money/testnet) contains latest configuration files for the testnet.

## Production Environment

**NOTE**: This guide only covers general settings for a production-level full node. You can find further details on considerations for operating a validator node in our [Validator Guide](https://docs.terra.money/validator/)

For the moment, this guide has only been tested against RPM-based Linux distributions.

### Increase Maximum Open Files

`terrad` can open more than 1024 files (which is default maximum) concurrently.
You will want to increase this limit.

Modify `/etc/security/limits.conf` to raise the `nofile` capability.

```
*                soft    nofile          65535
*                hard    nofile          65535
```

### Create a Dedicated User

`terrad` does not require the super user account. We **strongly** recommend using a normal user to run `terrad`. However, during the setup process you'll need super user permission to create and modify some files.

### Firewall Configuration

`terrad` uses several TCP ports for different purposes.

- `26656` is the default port for the P2P protocol. This port is opened in order to communicate with other nodes, and must be open to join a network. **However,** it does not have to be open to the public. For validator nodes, we recommend configuring `persistent_peers` and closing this port to the public.

- `26657` is the default port for the RPC protocol. This port is used for querying / sending transactions. In other words, this port needs to be opened for serving queries from `terracli`. It is safe to _NOT_ to open this port to the public unless you are planning to run a public node.

- `1317` is the default port for [Lite Client Daemon](https://docs.terra.money/terracli/lcd.html) (LCD), which can be enabled at `~/.terra/config/app.toml`. LCD provides HTTP RESTful API layer to allow applications and services to interact with your `terrad` instance through RPC. Check the [Terra REST API](https://swagger.terra.money) for usage examples. You don't need to open this port unless you have use of it.

- `26660` is the default port for interacting with the [Prometheus](https://prometheus.io) database which can be used for monitoring the environment. This port is not opened in the default configuration.

### Running Server as a Daemon

It is important to keep `terrad` running at all times. There are several ways to achieve this, and the simplest solution we recommend is to register `terrad` as a `systemd` service so that it will automatically get started upon system reboots and other events.

### Register terrad as a service

First, create a service definition file in `/etc/systemd/system`.

#### Sample file: `/etc/systemd/system/terrad.service`

```
[Unit]
Description=Terra Daemon
After=network.target

[Service]
Type=simple
User=terra
ExecStart=/data/terra/go/bin/terrad start
Restart=on-abort

[Install]
WantedBy=multi-user.target

[Service]
LimitNOFILE=65535
```

Modify the `Service` section from the given sample above to suit your settings.
Note that even if we raised the number of open files for a process, we still need to include `LimitNOFILE`.

After creating a service definition file, you should execute `systemctl daemon-reload`.

### Controlling the service

Use `systemctl` to control (start, stop, restart)

```bash
# Start
systemctl start terrad
# Stop
systemctl stop terrad
# Restart
systemctl restart terrad
```

### Accessing logs

```bash
# Entire log
journalctl -t terrad
# Entire log reversed
journalctl -t terrad -r
# Latest and continuous
journalctl -t terrad -f
```

## Resources

- Developer Tools

  - SDKs
    - [Terra.js](https://www.github.com/terra-money/terra.js) for JavaScript
    - [Jigu](https://www.github.com/terra-money/jigu) for Python
  - [Faucet](https://faucet.terra.money) can be used to get tokens for testnets
  - [LocalTerra](https://www.github.com/terra-money/LocalTerra) can be used to set up a private local testnet with configurable world state

- Block Explorers

  - [Terra Finder](https://finder.terra.money)
  - [Figment Hubble](https://hubble.figment.network/terra/chains/columbus-3)
  - [Stake ID by StakingFund](https://terra.stake.id)

- Wallets

  - [Terra Station](https://station.terra.money)
  - [Lunie](https://lunie.io/)

- Research

  - [Agora](https://agora.terra.money) - Research forum
  - [White Paper](https://terra.money/static/Terra_White_Paper.pdf)

## Community

- [Offical Website](https://terra.money)
- [Discord](https://discord.gg/Gutqybc)
- [Telegram](https://t.me/terra_announcements)
- [Twitter](https://twitter.com/terra_money)
- [YouTube](https://goo.gl/3G4T1z)

## Contributing

We are currently finalizing contribution standards and guidelines. In the meanwhile, if you are interested in contributing to the Terra Project, please contact our [admin](mailto:core@terra.money).

## License

This software is licensed under the Apache 2.0 license. Read more about it [here](LICENSE).

© 2020 Terraform Labs, PTE LTD

<hr/>

<p>&nbsp;</p>
<p align="center">
    <a href="https://terra.money/"><img src="http://terra.money/logos/terra_logo.svg" align="center" width=200/></a>
</p>
<div align="center">
  <sub><em>Powering the innovation of money.</em></sub>
</div>
