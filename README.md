## What is Terra Classic?

**[Terra](https://terra.money)** is a public, open-source blockchain protocol that provides fundamental infrastructure for a decentralized economy and enables open participation in the creation of new financial primitives to power the innovation of money.


**Classic** is the reference implementation of the Terra protocol, written in Golang. Terra Core is built atop [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and uses [Tendermint](https://github.com/tendermint/tendermint) BFT consensus. If you intend to work on Terra Core source, it is recommended that you familiarize yourself with the concepts in those projects.

Upon the implosion of Terra, a group of rebels seized control of the blockchain.  Terra's future is uncertain, but the rebels are now firmly in control. 

## Installation

### Binaries

The easiest way to get started is by downloading a pre-built binary for your operating system. You can find the latest binaries on the [releases](https://github.com/classic-terra/core/releases) page.

### From Source

**Step 1. Install Golang**

Go v1.20 is required for Terra Core.

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

**Step 2: Get Terra Core source code**

Use `git` to retrieve Terra Core from the [official repo](https://github.com/terra-money/core/) and checkout the `main` branch. This branch contains the latest stable release, which will install the `terrad` binary.

```bash
git clone https://github.com/classic-terra/core/
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
version: <version>
commit: <git commit hash>
..
...
....
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

### Join a testnet

Several testnets might exist simultaneously. Ensure that your version of `terrad` is compatible with the network you want to join.

To set up a node on the latest testnet, visit [the testnet repo](https://github.com/terra-money/testnet).

### Run a local testnet

The easiest way to set up a local testing environment is to run [LocalTerra](https://github.com/terra-money/LocalTerra), which automatically orchestrates a complete testing environment suited for development with zero configuration.

### Run a single node testnet

You can also run a local testnet using a single node. On a local testnet, you will be the sole validator signing blocks.


**Step 1. Create network and account**

First, initialize your genesis file to bootstrap your network. Create a name for your local testnet and provide a moniker to refer to your node:

```bash
terrad init --chain-id=<testnet_name> <node_moniker>
```

Next, create a Terra account by running the following command:

```bash
terrad keys add <account_name>
```

**Step 2. Add account to genesis**

Next, add your account to genesis and set an initial balance to start. Run the following commands to add your account and set the initial balance:

```bash
terrad add-genesis-account $(terrad keys show <account_name> -a) 100000000uluna,1000usd
terrad gentx <account_name> 10000000uluna --chain-id=<testnet_name>
terrad collect-gentxs
```

**Step 3. Run Terra daemon**

Now you can start your private Terra network:

```bash
terrad start
```

Your `terrad` node will be running a node on `tcp://localhost:26656`, listening for incoming transactions and signing blocks.

Congratulations, you've successfully set up your local Terra network!

## Set up a production environment

**NOTE**: This guide only covers general settings for a production-level full node. You can find further details on considerations for operating a validator node by visiting the [Terra validator guide](https://docs.terra.money/docs/full-node/manage-a-terra-validator/README.html).

This guide has been tested against Linux distributions only. To ensure you successfully set up your production environment, consider setting it up on an Linux system.

### Increase maximum open files

`terrad` can't open more than 1024 files (the default maximum) concurrently.

You can increase this limit by modifying `/etc/security/limits.conf` and raising the `nofile` capability.

```
*                soft    nofile          65535
*                hard    nofile          65535
```

### Create a dedicated user

It is recommended that you run `terrad` as a normal user. Super-user accounts are only recommended during setup to create and modify files.

### Port configuration

`terrad` uses several TCP ports for different purposes.

- `26656`: The default port for the P2P protocol. Use this port to communicate with other nodes. While this port must be open to join a network, it does not have to be open to the public. Validator nodes should configure `persistent_peers` and close this port to the public.

- `26657`: The default port for the RPC protocol. This port is used for querying / sending transactions and must be open to serve queries from `terrad`. **DO NOT** open this port to the public unless you are planning to run a public node.

- `1317`: The default port for [Lite Client Daemon](https://docs.terra.money/docs/develop/how-to/start-lcd.html) (LCD), which can be enabled in `~/.terra/config/app.toml`. The LCD provides an HTTP RESTful API layer to allow applications and services to interact with your `terrad` instance through RPC. Check the [Terra REST API](https://lcd.terra.dev/swagger/#/) for usage examples. Don't open this port unless you need to use the LCD.

- `26660`: The default port for interacting with the [Prometheus](https://prometheus.io) database. You can use Promethues to monitor an environment. This port is closed by default.

### Run the server as a daemon

**Important**:

Keep `terrad` running at all times. The simplest solution is to register `terrad` as a `systemd` service so that it automatically starts after system reboots and other events.


### Register terrad as a service

First, create a service definition file in `/etc/systemd/system`.

**Sample file: `/etc/systemd/system/terrad.service`**

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
Note that even if you raised the number of open files for a process, you still need to include `LimitNOFILE`.

After creating a service definition file, you should execute `systemctl daemon-reload`.

### Start, stop, or restart service

Use `systemctl` to control (start, stop, restart)

```bash
# Start
systemctl start terrad
# Stop
systemctl stop terrad
# Restart
systemctl restart terrad
```

### Access logs

```bash
# Entire log
journalctl -t terrad
# Entire log reversed
journalctl -t terrad -r
# Latest and continuous
journalctl -t terrad -f
```

## Using `docker-compose`

1. Install Docker

	- [Docker Install documentation](https://docs.docker.com/install/)
	- [Docker-Compose Install documentation](https://docs.docker.com/compose/install/)

2. Create a new folder on your local machine and copy docker-compose\docker-compose.yml

3. Review the docker-compose.yml contents

4. Bring up your stack by running

	```bash
	docker-compose up -d
	```

5. Add your wallet
    ```bash
	docker-compose exec node sh /keys-add.sh
	```

6. Copy your terra wallet address and go to the terra faucet here -> http://45.79.139.229:3000/ Put your address in and give yourself luna coins.

7. Start the validator
	```bash
	docker-compose exec node sh /create-validator.sh
	```

### Cheat Sheet:

#### Start

```bash
docker-compose up -d
```

#### Stop

```bash
docker-compose down
```

#### View Logs

```bash
docker-compose logs -f
```

#### Run Terrad Commands Example

```bash
docker-compose exec node terrad status
```

#### Upgrade

```bash
docker-compose down
docker-compose pull
docker-compose up -d
```

#### Build from source

```sh
make build-all -f contrib/terra-operator/Makefile
```

## Resources

- Developer Tools

  - Terra developer documentation(https://docs.terra.money)
  - [TerraWiki.org](https://terrawiki.org) - The Terra community wiki.
  - SDKs
    - [Terra.js](https://www.github.com/terra-money/terra.js) for JavaScript
    - [terra-sdk-python](https://www.github.com/terra-money/terra-sdk-python) for Python
  - [Faucet](https://faucet.terra.money) can be used to get tokens for testnets
  - [LocalTerra](https://www.github.com/terra-money/LocalTerra) can be used to set up a private local testnet with configurable world state

- Developer Forums
  - [Terra Developer Discord](https://discord.com/channels/464241079042965516/591812948867940362)
  - [Terra DEveloper Telegram room](https://t.me/+gCxCPohmVBkyNDRl)


- Block Explorers

  - [Terra Finder](https://finder.terra.money) - Terra's basic block explorer.
  - [Terrascope](https://terrascope.info/) - A community-run block explorer with extra features.
  - [Stake ID](https://terra.stake.id) - A block explorer made by Staking Fund
  - [Hubble](https://hubble.figment.network/terra/chains/columbus-5) - by Figment

- Wallets

  - [Terra Station](https://station.terra.money) - The official Terra wallet.
  - Terra Station Mobile
    - [iOS](https://apps.apple.com/us/app/terra-station/id1548434735)
    - [Android](https://play.google.com/store/apps/details?id=money.terra.station&hl=en_US&gl=US)
    
  - [Falcon Wallet](https://falconwallet.app/)
  - [Leap Wallet](https://chrome.google.com/webstore/detail/leap-wallet/aijcbedoijmgnlmjeegjaglmepbmpkpi/?utm_source=Leap&utm_medium=Bio&utm_campaign=Leap)
  - [XDeFi](https://chrome.google.com/webstore/detail/xdefi-wallet/hmeobnfnfcmdkdcmlblgagmfpfboieaf)
  - [Liquality](https://liquality.io/)

- Research

  - [Agora](https://agora.terra.money) - Research forum
  - [White Paper](https://assets.website-files.com/611153e7af981472d8da199c/618b02d13e938ae1f8ad1e45_Terra_White_paper.pdf)

## Community

- [Offical Website](https://terra.money)
- [Discord](https://discord.gg/e29HWwC2Mz)
- [Telegram](https://t.me/terra_announcements)
- [Twitter](https://twitter.com/terra_money)
- [YouTube](https://goo.gl/3G4T1z)

## Contributing

If you are interested in contributing to Terra Core source, please review our [code of conduct](./CODE_OF_CONDUCT.md).

## License

This software is licensed under the Apache 2.0 license. Read more about it [here](LICENSE).

© 2022 Terraform Labs, PTE LTD

<hr/>

<p>&nbsp;</p>
<p align="center">
    <a href="https://terra.money/"><img src="https://assets.website-files.com/611153e7af981472d8da199c/61794f2b6b1c7a1cb9444489_symbol-terra-blue.svg" align="center" width=200/></a>
</p>
<div align="center">
  <sub><em>Powering the innovation of money.</em></sub>
</div>
