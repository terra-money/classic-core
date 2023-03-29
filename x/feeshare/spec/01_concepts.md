<!--
order: 1
-->

# Concepts

## FeeShare

The FeeShare module is a revenue-per-transaction model, which allows developers to get paid for deploying their decentralized applications (dApps) on Terra Classic. This helps developers to generate revenue every time a user interacts with their contracts on the chain. This registration is permissionless to sign up for and begin earning fees from. By default, 50% of all transaction fees for Execute Messages are shared. This can be changed by governance and implemented by the `x/feeshare` module.

## Registration

Developers register their contract applications to gain their cut of fees per execution. Any contract can be registered by a developer by submitting a signed transaction. After the transaction is executed successfully, the developer will start receiving a portion of the transaction fees paid when a user interacts with the registered contract. The developer can have the funds sent to their wallet, a DAO, or any other wallet address on the Terra Classic network.

::: tip
 **NOTE**: If your contract is part of a development project, please ensure that the deployer of the contract (or the factory/DAO that deployed the contract) is an account that is owned by that project. This avoids the situation, that an individual deployer who leaves your project could become malicious.
:::

## Fee Distribution

As described above, developers will earn a portion of the transaction fee after registering their contracts. To understand how transaction fees are distributed, we will look at the following in detail:

* The transactions eligible are only [Wasm Execute Txs](https://github.com/CosmWasm/wasmd/blob/main/proto/cosmwasm/wasm/v1/tx.proto#L115-L127) (`MsgExecuteContract`).

### WASM Transaction Fees

Users pay transaction fees to pay to interact with smart contracts on Terra Classic. When a transaction is executed, the entire fee amount (`gas limit * gas price`) is sent to the `FeeCollector` module account during the [Cosmos SDK AnteHandler](https://docs.cosmos.network/main/modules/auth/#antehandlers) execution. After this step, the `FeeCollector` sends 50% of the funds and splits them between contracts that were executed on the transaction. If the fees paid are not accepted by governance, there is no payout to the developers (for example, niche base tokens) for tax purposes. If a user sends a message and it does not interact with any contracts (ex: bankSend), then the entire fee is sent to the `FeeCollector` as expected.
