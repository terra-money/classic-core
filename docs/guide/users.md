# For non-validating users

## Terra Accounts

At the core of every Terra account, there is a seed, which takes the form of a 12 or 24-words mnemonic. From this mnemonic, it is possible to create any number of Terra accounts, i.e. pairs of private key/public key. This is called an HD wallet \(see [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) for more information on the HD wallet specification\).

```text
     Account 0                         Account 1                         Account 2

+------------------+              +------------------+               +------------------+
|                  |              |                  |               |                  |
|    Address 0     |              |    Address 1     |               |    Address 2     |
|        ^         |              |        ^         |               |        ^         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        +         |              |        +         |               |        +         |
|  Public key 0    |              |  Public key 1    |               |  Public key 2    |
|        ^         |              |        ^         |               |        ^         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        +         |              |        +         |               |        +         |
|  Private key 0   |              |  Private key 1   |               |  Private key 2   |
|        ^         |              |        ^         |               |        ^         |
+------------------+              +------------------+               +------------------+
         |                                 |                                  |
         |                                 |                                  |
         |                                 |                                  |
         +--------------------------------------------------------------------+
                                           |
                                           |
                                 +---------+---------+
                                 |                   |
                                 |  Mnemonic (Seed)  |
                                 |                   |
                                 +-------------------+
```

The funds stored in an account are controlled by the private key. This private key is generated using a one-way function from the mnemonic. If you lose the private key, you can retrieve it using the mnemonic. However, if you lose the mnemonic, you will lose access to all the derived private keys. Likewise, if someone gains access to your mnemonic, they gain access to all the associated accounts.

{% hint style="danger" %}
Do not lose or share your 24 words with anyone.

To prevent theft or loss of funds, it is best to ensure that you keep multiple copies of your mnemonic, and store it in a safe, secure place and that only you know how to access.

If someone is able to gain access to your mnemonic, they will be able to gain access to your private keys and control the accounts associated with them.
{% endhint %}

The address is a public string with a human-readable prefix \(e.g. `terra10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg`\) that identifies your account. When someone wants to send you funds, they send it to your address. It is computationally infeasible to find the private key associated with a given address.

### Redeem Luna from the sale

If you participated in the fundraiser, you should be in possession of a 24-word mnemonic created through the Registrar tool. If you have not yet generated your wallet, please contact the Terra team right away.

#### On a ledger device

At the core of a ledger device, there is a mnemonic used to generate accounts on multiple blockchains \(including Terra\). Usually, you will create a new mnemonic when you initialize your ledger device. However, it is possible to tell the ledger device to use a mnemonic provided by the user instead. Let us go ahead and see how you can input the mnemonic you obtained during the fundraiser as the seed of your ledger device.

{% hint style="warning" %}
To do this, **it is preferable to use a brand new ledger device** as there can be only one mnemonic per ledger device.

If, however, you want to use a ledger that is already initialized with a seed, you can reset it by going in `Settings`&gt;`Device`&gt;`Reset All`.

**Please note that this will wipe out the seed currently stored on the device.**

**If you have not properly secured the associated mnemonic, you could lose your funds!!!**
{% endhint %}

The following steps need to be performed on an un-initialized ledger device:

1. Connect your ledger device to the computer via USB
2. Press both buttons
3. Do **NOT** choose the "Config as a new device" option. Instead, choose "Restore Configuration"
4. Choose a PIN
5. Choose the 24 words option
6. Input each of the words you got during the fundraiser, in the correct order. 

Your ledger is now correctly set up with your fundraiser mnemonic! Do not lose this mnemonic! If your ledger is compromised, you can always restore a new device again using the same mnemonic.

Next, click [here](users.md#using-a-ledger-device) to learn how to generate an account.

#### On a computer

{% hint style="warning" %}
It is more secure to perform this action on an offline computer.
{% endhint %}

To restore an account using a fundraiser mnemonic and store the associated encrypted private key on a computer, use the following command:

```bash
terracli keys add <yourKeyName> --recover
```

You will be prompted to input a passphrase that is used to encrypt the private key of account `0` on disk. Each time you want to send a transaction, this password will be required. If you lose the password, you can always recover the private key with the mnemonic.

* `<yourKeyName>` is the name of the account. It is a reference to the account number used to derive the key pair from the mnemonic. You will use this name to identify your account when you want to send a transaction.
* You can add the optional `--account` flag to specify the path \(`0`, `1`, `2`, ...\) you want to use to generate your account. By default, account `0` is generated. 

### Creating an account

To create an account, you just need to have `terracli` installed. Before creating it, you need to know where you intend to store and interract with your private keys. The best options are to store them in an offline dedicated computer or a ledger device. Storing them on your regular online computer involves more risk, since anyone who infiltrates your computer through the internet could exfiltrate your private keys and steal your funds.

#### Using a ledger device

{% hint style="warning" %}
Only use Ledger devices that you bought factory new or trust fully.
{% endhint %}

When you initialize your ledger, a 24-word mnemonic is generated and stored in the device. This mnemonic is compatible with Terra and Terra accounts can be derived from it. Therefore, all you have to do is make your ledger compatible with `terracli`. To do so, you need to go through the following steps:

1. Download the Ledger Live app [here](https://www.ledger.com/pages/ledger-live). 
2. Connect your ledger via USB and update to the latest firmware
3. Go to the ledger live app store, and download the "Terra" application \(this can take a while\). **Note: You may have to enable** `Dev Mode` **in the** `Settings` **of Ledger Live to be able to download the "Terra" application**.
4. Navigate to the Terra app on your ledger device

Then, to create an account, use the following command:

```bash
terracli keys add <yourAccountName> --ledger
```

{% hint style="warning" %}
This command will only work while the Ledger is plugged in and unlocked.
{% endhint %}

* `<yourKeyName>` is the name of the account. It is a reference to the account number used to derive the key pair from the mnemonic. You will use this name to identify your account when you want to send a transaction.
* You can add the optional `--account` flag to specify the path \(`0`, `1`, `2`, ...\) you want to use to generate your account. By default, account `0` is generated. 

#### Using a computer

{% hint style="warning" %}
It is more secure to perform this action on an offline computer.
{% endhint %}

To generate an account, just use the following command:

```bash
terracli keys add <yourKeyName>
```

The command will generate a 24-words mnemonic and save the private and public keys for account `0` at the same time. You will be prompted to input a passphrase that is used to encrypt the private key of account `0` on disk. Each time you want to send a transaction, this password will be required. If you lose the password, you can always recover the private key with the mnemonic.

{% hint style="danger" %}
Do not lose or share your 24 words with anyone.

To prevent theft or loss of funds, it is best to ensure that you keep multiple copies of your mnemonic, and store it in a safe, secure place and that only you know how to access.

If someone is able to gain access to your mnemonic, they will be able to gain access to your private keys and control the accounts associated with them.
{% endhint %}

{% hint style="info" %}
 After you have secured your mnemonic \(triple check!\), you can delete bash history to ensure no one can retrieve it.

```bash
history -c
rm ~/.bash_history
```
{% endhint %}

* `<yourKeyName>` is the name of the account. It is a reference to the account number used to derive the key pair from the mnemonic. You will use this name to identify your account when you want to send a transaction.
* You can add the optional `--account` flag to specify the path \(`0`, `1`, `2`, ...\) you want to use to generate your account. By default, account `0` is generated. 

You can generate more accounts from the same mnemonic using the following command:

```bash
terracli keys add <yourKeyName> --recover --account 1
```

This command will prompt you to input a passphrase as well as your mnemonic. Change the account number to generate a different account.

## Accessing the Terra network

In order to query the state and send transactions, you need a way to access the network. To do so, you can either run your own full-node, or connect to someone else's.

{% hint style="danger" %}
Do not share your mnemonic (12 or 24 words) with anyone.

The only person who should ever need to know it is you.

This is especially important if you are ever approached via email or direct message by someone requesting that you share your mnemonic for any kind of blockchain services or support.

No one from Terra will ever send an email that asks for you to share any kind of account credentials or your mnemonic.
{% endhint %}

### Running your own full-node

This is the most secure option, but comes with relatively high resource requirements. In order to run your own full-node, you need good bandwidth and at least 1TB of disk space.

You will find the tutorial on how to install `terrad` [here](installation.md), and the guide to run a full-node [here](join-network.md).

### Connecting to a remote full-node

If you do not want or cannot run your own node, you can connect to someone else's full-node. You should pick an operator you trust, because a malicious operator could return incorrect query results or censor your transactions. However, they will never be able to steal your funds, as your private keys are stored locally on your computer or ledger device. Possible options of full-node operators include validators, wallet providers or exchanges.

In order to connect to the full-node, you will need an address of the following form: `https://77.87.106.33:26657` \(_Note: This is a placeholder_\). This address has to be communicated by the full-node operator you choose to trust. You will use this address in the [following section](users.md#setting-up-terracli).

## Setting up `terracli`

{% hint style="warning" %}
Please check that you are always using the latest stable release of `terracli`.
{% endhint %}

`terracli` is the tool that enables you to interact with the node that runs on the Terra Protocol network, whether you run it yourself or not. Let us set it up properly.

In order to set up `terracli`, use the following command:

```bash
terracli config <flag> <value>
```

It allows you to set a default value for each given flag.

First, set up the address of the full-node you want to connect to:

```bash
terracli config node <host>:<port

// example: terracli config node https://77.87.106.33:26657
```

If you run your own full-node, just use `tcp://localhost:26657` as the address.

Then, let us set the default value of the `--trust-node` flag:

```bash
terracli config trust-node false

// Set to true if you run a light-client node, false otherwise
```

Finally, let us set the `chain-id` of the blockchain we want to interact with:

```bash
terracli config chain-id gos-6
```

## Querying the state

{% hint style="warning" %}
Before you can bond luna and withdraw rewards, you need to [set up](users.md#setting-up-terracli) `terracli`.
{% endhint %}

`terracli` lets you query all relevant information from the blockchain, like account balances, amount of bonded tokens, outstanding rewards, and more. Next is a list of the most useful commands for delegator.

```bash
// query account balances and other account-related information
terracli query account

// query the list of validators
terracli query staking validators

// query the information of a validator given their address (e.g. terravaloper1n5pepvmgsfd3p2tqqgvt505jvymmstf6s9gw27)
terracli query staking validator <validatorAddress>

// query all delegations made from a delegator given their address (e.g. terra10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg)
terracli query staking delegations <delegatorAddress>

// query a specific delegation made from a delegator (e.g. terra10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg) to a validator (e.g. terravaloper1n5pepvmgsfd3p2tqqgvt505jvymmstf6s9gw27) given their addresses
terracli query staking delegation <delegatorAddress> <validatorAddress>

// query the rewards of a delegator given a delegator address (e.g. terra10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg)
terracli query distr rewards <delegatorAddress>
```

For more commands, just type:

```bash
terracli query
```

For each command, you can use the `-h` or `--help` flag to get more information.

## Sending Transactions

### A note on gas and fees

Transactions on the Terra Protocol network need to include a transaction fee in order to be processed. This fee pays for the gas required to run the transaction. The formula is the following:

```text
fees = gas * gasPrices
```

The `gas` is dependent on the transaction. Different transaction require different amount of `gas`. The `gas` amount for a transaction is calculated as it is being processed, but there is a way to estimate it beforehand by using the `auto` value for the `gas` flag. Of course, this only gives an estimate. You can adjust this estimate with the flag `--gas-adjustment` \(default `1.0`\) if you want to be sure you provide enough `gas` for the transaction.

The `gasPrice` is the price of each unit of `gas`. Each validator sets a `min-gas-price` value, and will only include transactions that have a `gasPrice` greater than their `min-gas-price`.

The transaction `fees` are the product of `gas` and `gasPrice`. As a user, you have to input 2 out of 3. The higher the `gasPrice`/`fees`, the higher the chance that your transaction will get included in a block.
