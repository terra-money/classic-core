# Join a network

{% hint style="info" %}
See the [testnet repo](https://github.com/terra-project/networks) for information on the latest testnet, including the correct version of the Terra Core to use and details about the genesis file.
{% endhint %}

{% hint style="warning" %}
You need to [install terra](installation.md) before going further.
{% endhint %}

## Setting Up a New Node

{% hint style="info" %}
If you ran a full node on a previous testnet, please skip to [Upgrading From Previous Testnet](join-network.md#upgrading-from-previous-testnet).
{% endhint %}

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```bash
terrad init <your_custom_moniker>
```

{% hint style="warning" %}
Monikers can only contain  ASCII characters.
Using Unicode characters will render your node unreachable.
{% endhint %}

You can edit this `moniker` later, in the `~/.terrad/config/config.toml` file:

```text
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```

You can edit the `~/.terrad/config/terrad.toml` file in order to enable the anti spam mechanism and reject incoming transactions with less than a minimum fee:

```text
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# Validators reject any tx from the mempool with less than the minimum fee per gas.
minimum_fees = ""
```

Your full node has been initialized! Please skip to [Genesis & Seeds](join-network.md#genesis-and-seeds).

## Upgrading From Previous Testnet

These instructions are for full nodes that have ran on previous testnets and would like to upgrade to the latest testnet.

### Reset Data

First, remove the outdated files and reset the data.

```bash
rm $HOME/.terrad/config/addrbook.json $HOME/.terrad/config/genesis.json
terrad unsafe-reset-all
```

Your node is now in a pristine state while keeping the original `priv_validator.json` and `config.toml`. If you had any sentry nodes or full nodes setup before, your node will still try to connect to them, but may fail if they haven't also been upgraded.

{% hint style="danger" %}
Make sure that every node has a unique `priv_validator.json`.

Do not copy the `priv_validator.json` from an old node to multiple new nodes.
Running two nodes with the same `priv_validator.json` will cause you to double sign.
{% endhint %}

### Software Upgrade

Now it is time to upgrade the software. Go to the project directory, and run:

```bash
git checkout master && git pull
make
```

{% hint style="info" %}
If you have issues at this step, please check that you have the latest stable version of GO installed.
{% endhint %}

Note we use `master` here since it contains the latest stable release. See the [testnet repo](https://github.com/terra-project/networks) for details on which version is needed for which testnet, and the [SDK release page](https://github.com/terra-project/core//releases) for details on each release.

Your full node has been cleanly upgraded!

## Genesis & Seeds

### Copy the Genesis File

Fetch the testnet's `genesis.json` file into `terrad`'s config directory.

```bash
mkdir -p $HOME/.terrad/config
curl https://raw.githubusercontent.com/terra-project/launch/master/genesis.json > $HOME/.terrad/config/genesis.json
```

Note we use the `latest` directory in the [networks repo](https://github.com/terra-project/networks) which contains details for the latest testnet. If you are connecting to a different testnet, ensure you get the right files.

To verify the correctness of the configuration run:

```bash
terrad start
```

### Add Seed Nodes

Your node needs to know how to find peers. You'll need to add healthy seed nodes to `$HOME/.terrad/config/config.toml`. The `testnets` repo contains links to the seed nodes for each testnet. If you are looking to join the running testnet please [check the repository for details](https://github.com/terra-project/networks) on which nodes to use.

If those seeds aren't working, you can find more seeds and persistent peers on the [Terra Station](https://station.terra.money). Open the the `Full Nodes` pane and select nodes that do not have private \(`10.x.x.x`\) or [local IP addresses](https://en.wikipedia.org/wiki/Private_network). The `Persistent Peer` field contains the connection string. For best results use 4-6.

For more information on seeds and peers, you can [read this](https://github.com/tendermint/tendermint/blob/develop/docs/tendermint-core/using-tendermint.md#peers).

## Run a Full Node

Start the full node with this command:

```bash
terrad start
```

Check that everything is running smoothly:

```bash
terracli status
```

View the status of the network with the [Terra Finder](https://finder.terra.money). Once your full node syncs up to the current block height, you should see it appear on the [list of full nodes](https://terra.stake.id/).

## Export State

Terra can dump the entire application state to a JSON file, which could be useful for manual analysis and can also be used as the genesis file of a new network.

Export state with:

```bash
terrad export > [filename].json
```

You can also export state from a particular height \(at the end of processing the block of that height\):

```bash
terrad export --height [height] > [filename].json
```

If you plan to start a new network from the exported state, export with the `--for-zero-height` flag:

```bash
terrad export --height [height] --for-zero-height > [filename].json
```

## Upgrade to Validator Node

You now have an active full node. What's the next step? You can upgrade your full node to become a Terra Validator. The top 100 validators have the ability to propose new blocks to the Terra network. Continue onto [the Validator Setup](validators.md).
