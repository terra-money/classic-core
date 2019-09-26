# Ledger Nano support for Terra

A Ledger Nano S is a hardware wallet. Hardware wallets are considered very secure for the storage of a user’s private keys in the blockchain world. Your digital assets are safe even when using an infected (or untrusted) PC. Follow these instructions to interact with the Terra blockchain using a Ledger Nano device. 

## Requirements
- You have [initialized your Ledger Nano S](https://support.ledgerwallet.com/hc/en-us/articles/360000613793)
- The latest firmware is [installed](https://support.ledger.com/hc/en-us/articles/360002731113-Update-Ledger-Nano-S-firmware)
- Ledger Live is [ready to use](https://support.ledger.com/hc/en-us/articles/360006395233-Take-your-first-steps)
- Google Chrome is installed.

## Installation

{% hint style="warning" %}
To do this, it is preferable to use a brand new ledger device as there can be only one mnemonic per ledger device.

If, however, you want to use a ledger that is already initialized with a seed, you can reset it by going in `Settings`>`Device`>`Reset All`.

**Please note that this will wipe out the seed currently stored on the device.**

**If you have not properly secured the associated mnemonic, you could lose your funds!**
{% endhint %}


- Open the **Manager** in Ledger Live.
- Connect and unlock your Ledger Nano S.
- If asked, allow the manager on your device by pressing the right button.
- Find **Terra** in the app catalog.
- Click the Install button of the app.
- An installation window appears.
- Your device will display **Processing…**
- The app installation is confirmed.

## Setup 

Before we can configure the Ledger Nano S to interact with the Terra blockchain, we need the following: 

- [Generating new keys or recovering fundraiser keys from a ledger device](./users.md#on-a-ledger-device)
- [Creating an account using a ledger device](./users.md#using-a-ledger-device)
- [A running `terrad` instance connected to the network you wish to use.](./users.md#accessing-the-terra-network)
- [A `terracli` instance configured to connect to your chosen `terrad` instance.](./users.md#setting-up-terracli)

Now, you are all set to start sending and receiving transactions on the network.

## Use the Ledger with the CLI

{% hint style="info" %}
Your ledger device must be on and the Terra ledger app must be in the foreground to perform the following actions. 
{% endhint %}

### How to view account balance

You can [use `terracli` to view the account balance](./terracli.md####Query-Account-balance), using the key created in the above step. 

### How to receive tokens

Run `terracli keys show <yourAccountName> -d` to see the account address at which to receive tokens. You will be asked to confirm that the intended receiving address (matching `<yourAccountName>`) matches the account registered on the ledger. If the two addresses match, press the right button to confirm, and offer the address to the sender. If they do not match, press the left button to reject. 

### How to send tokens

1. You can [use `terracli` to send tokens](./terracli.md###Send-Tokens), using the key created in the above step. 
2. You will be asked to confirm the details of the transaction. Before confirming, check that the destination address on the Ledger display matches your intended destination address. 

## Support

Please speak to us on our [public discord group](https://discord.gg/a4tqbt) to drop us a line if you run into problems in setup or usage of the Ledger device. 
