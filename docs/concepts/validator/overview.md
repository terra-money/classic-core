# Validators Overview

## Introduction

Terra Core is based on Tendermint, which relies on a set of validators that are responsible for committing new blocks in the blockchain. These validators participate in the consensus protocol by broadcasting votes which contain cryptographic signatures signed by each validator's private key.

Validator candidates can bond their own Luna and have Luna "delegated", or staked, to them by token holders. The Columbus Mainnet will have 100 validators, but over time this will increase to 300 validators according to a predefined schedule. The validators are determined by who has the most stake delegated to them — the top 100 validator candidates with the most stake will become Terra validators.

Validators and their delegators will earn the following fees: 

- **Compute fees**: To prevent spamming, validators may set minimum gas fees for transactions to be included in their mempool. At the end of every block, the compute fees are disbursed to the participating validators pro-rata to stake. 
- **Stability fees**: To stabilize the value of Luna, the protocol charges a small percentage transaction fee ranging from 0.1% to 1% on every Terra transaction, capped at 1 TerraSDR. This is paid in any Terra currency, and is disbursed pro-rata to stake at the end of every block in TerraSDR. 
- **Seigniorage rewards**: To stabilize the value of Luna, the protocol commits to using some variable portion of Terra seigniorage (see the market and treasury modules for how this functions) to buy back and burn Luna tokens. This creates scarcity for Luna tokens and indirectly rewards validators. 

Note that validators can set commission on the fees their delegators receive as additional incentive.

If validators double sign, are frequently offline or do not participate in governance, their staked Luna (including Luna of users that delegated to them) can be slashed. The penalty depends on the severity of the violation.

## Hardware

There currently exists no appropriate cloud solution for validator key management. This may change in 2018 when cloud SGX becomes more widely available. For this reason, validators must set up a physical operation secured with restricted access. A good starting place, for example, would be co-locating in secure data centers.

Validators should expect to equip their datacenter location with redundant power, connectivity, and storage backups. Expect to have several redundant networking boxes for fiber, firewall and switching and then small servers with redundant hard drive and failover. Hardware can be on the low end of datacenter gear to start out with.

We anticipate that network requirements will be low initially. The current testnet requires minimal resources. Then bandwidth, CPU and memory requirements will rise as the network grows. Large hard drives are recommended for storing years of blockchain history.

## Set Up a Website

Set up a dedicated validator's website and signal your intention to become a validator by contacting Terraform Labs (core at terra dot money). This is important since delegators will want to have information about the entity they are delegating their Luna to.

## Seek Legal Advice

Seek legal advice if you intend to run a Validator.
